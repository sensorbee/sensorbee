package exp

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func setUpRun() cli.Command {
	cmd := cli.Command{
		Name:   "run",
		Usage:  "run BQL statements with cache",
		Action: runRun,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "remove-cache",
				Usage: "remove outdated cache files after execution is done",
			},
			cli.BoolFlag{
				Name:  "ignore-cache",
				Usage: "run all statements again by ignoring existing cache files",
			},
		},
	}

	// TODO: flags
	// --ignore-cache: run all statements again.
	// --from <hash>: invalidate cache of the statement having the hash and run it again.
	// --cache <name> (default .sensorbee-exp-cache)
	// --leave-intermediate-cache: create a cache file not only on the last node but also for intermediate nodes
	// --limit <num>: the maximum number of tuples generated from the source.
	return cmd
}

// TODO: write a note that all sources mustn't be rewindable and must be paused.
// They have to be stopped after they emit all tuples.

// TODO: support save/load of UDSs
//       UDSs need to be saved on each step

// TODO: add subcommand which deletes all caches except the latest ones.
// TODO: support annotation like "-- @something" in BQL statements (only at the beginning of each line)

// TODO: reduce number of files to be generated

const (
	cacheInfoFilename = ".sensorbee-exp-cache" // TODO: this should be an option
)

func runRun(c *cli.Context) {
	err := func() error {
		if a := c.Args(); len(a) != 1 {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		}
		filename := c.Args()[0]

		body, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot read BQL statements from %v: %v", filename, err)
		}

		cache, err := LoadCacheFromFile(cacheInfoFilename)
		if err != nil {
			return fmt.Errorf("cannot load cache information:%v", err)
		}

		stmts, err := Parse(string(body))
		if err != nil {
			return fmt.Errorf("cannot parse BQL statements:\n%v", err)
		}

		ith := -1
		if !c.Bool("ignore-cache") {
			ith = findCache(cache, stmts)
		}
		cache.TruncateAfter(ith)
		ith++
		for i := ith; i < len(stmts.Stmts); i++ {
			fmt.Println("Run:", stmts.Stmts[i].String())
			if err := execute(cache, stmts, i); err != nil {
				return err
			}
			if err := SaveCacheToFile(cache, cacheInfoFilename); err != nil {
				return err
			}
		}
		if c.Bool("remove-cache") {
			if err := removeCacheFiles(cache.Outdated); err != nil {
				return err
			}
			cache.Outdated = nil
			if err := SaveCacheToFile(cache, cacheInfoFilename); err != nil {
				return err
			}
		}
		fmt.Println("Done")
		return nil
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func findCache(cache *Cache, stmts *Statements) int {
	if len(stmts.Stmts) == 0 {
		return -1
	}
	prevHash := cache.Seed
	for i, stmt := range stmts.Stmts {
		hash := hashStatement(prevHash, stmt)
		if !hasCache(cache, i, hash, stmt) {
			return i - 1
		}
		prevHash = hash
	}
	return len(stmts.Stmts) - 1
}

func hasCache(cache *Cache, ith int, hash string, stmt *Statement) bool {
	ent := cache.Get(ith)
	if ent == nil {
		return false
	}
	if ent.Hash != hash {
		return false
	}

	// Because STREAMs and INSERT INTOs may change the state of UDSs, all
	// states needs to be saved.
	if !(stmt.IsStream() || stmt.IsInsertStatement()) {
		// The statement doesn't have intermediate results.
		return true
	}

	// TODO: check cache information of all UDSs created so far

	if stmt.IsInsertStatement() {
		return true
	}

	_, err := os.Stat(CacheFilename(ent))
	if err != nil {
		return false
	}
	return true
}

func execute(cache *Cache, stmts *Statements, ith int) error {
	prevHash := ""
	if ith < 1 {
		prevHash = cache.Seed
	} else {
		prevHash = cache.Get(ith - 1).Hash
	}

	conf := &core.ContextConfig{}
	conf.Flags.DroppedTupleLog.Set(true)
	conf.Flags.DroppedTupleSummarization.Set(true)
	t := core.NewDefaultTopology(core.NewContext(conf), "sensorbee_exp")
	defer t.Stop()
	tb, err := bql.NewTopologyBuilder(t)
	if err != nil {
		return fmt.Errorf("cannot create a topology builder: %v", err)
	}

	curStmt := stmts.Stmts[ith]
	pausedSrcs, err := setUpTopology(tb, cache, stmts, ith)
	if err != nil {
		return err
	}

	resumeAll := func() error {
		for _, name := range pausedSrcs {
			s, err := tb.Topology().Source(name)
			if err != nil {
				return fmt.Errorf("cannot find the source %v: %v", name, err)
			}
			if err := s.Resume(); err != nil {
				return fmt.Errorf("cannot resume the source %v: %v", name, err)
			}
		}
		return nil
	}
	if n, err := tb.AddStmt(curStmt.Stmt); err != nil {
		return fmt.Errorf(`cannot add the statement "%v": %v`, curStmt, err)
	} else if n != nil && n.Type() == core.NTSource {
		sn, _ := n.(core.SourceNode)
		if _, ok := sn.Source().(core.RewindableSource); ok {
			return fmt.Errorf(`rewindable source "%v" isn't supported`, n.Name())
		}
		pausedSrcs = append(pausedSrcs, n.Name())
	}

	hash := hashStatement(prevHash, curStmt)
	ent := &CacheEntry{
		NodeName:  curStmt.NodeName(),
		Hash:      hash,
		Stmt:      curStmt.String(),
		Timestamp: time.Now(),
	}
	if err := executeTarget(t, ent, curStmt, resumeAll); err != nil {
		return err
	}

	// Because state can be changed through various statements, they have to
	// be saved on each state execution.
	if err := saveStates(t, ent); err != nil {
		return err
	}
	cache.Add(ent)
	return nil
}

func setUpTopology(tb *bql.TopologyBuilder, cache *Cache, stmts *Statements, ith int) ([]string, error) {
	curStmt := stmts.Stmts[ith]
	inputs, err := curStmt.Input()
	if err != nil {
		return nil, fmt.Errorf(`cannot get input information of the statement "%v": %v`, curStmt, err)
	}
	latestEnt := cache.Get(ith - 1) // this can be nil

	if len(inputs) > 1 {
		// TODO: reduce number of steps that execute goes back due to JOIN.
		// Executing from the nearest common ancestor might work fine.
		// Because INSERT INTO or other statements may change UDSs, it isn't
		// a simple problem.

		// TODO: add option to perform JOIN using cache without thinking processing lag.

		var pausedSrcs []string
		for i := 0; i < ith; i++ {
			stmt := stmts.Stmts[i]
			if stmt.IsInsertStatement() {
				// Although INSERT INTO might change UDSs, UDSs are cached on
				// each execution of a statement. Therefore, they don't have to
				// be run again. However, it means a UDF which modify UDSs can
				// result in a confusing behavior due to duplicated evaluation.
				continue
			}

			if stmt.IsCreateStateStatement() {
				if err := setUpState(tb, latestEnt, stmt); err != nil {
					return nil, err
				}
				continue
			}

			if n, err := addStmt(tb, stmt); err != nil {
				return nil, err
			} else if n != nil && n.Type() == core.NTSource {
				pausedSrcs = append(pausedSrcs, n.Name())
			}
		}
		return pausedSrcs, nil
	}

	for i := 0; i < ith; i++ {
		stmt := stmts.Stmts[i]
		if stmt.IsDataSourceNodeQuery() || stmt.IsInsertStatement() {
			continue
		}

		if stmt.IsCreateStateStatement() {
			if err := setUpState(tb, latestEnt, stmt); err != nil {
				return nil, err
			}
			continue
		}

		if _, err := addStmt(tb, stmt); err != nil {
			return nil, err
		}
	}

	// len(inputs) <= 1 due to the condition above.
	if len(inputs) == 0 {
		return nil, nil
	}

	ent := cache.Node(inputs[0])
	if ent == nil {
		return nil, fmt.Errorf("The node %v doesn't have cache.\nPlease run all without cache again.", inputs[0])
	}
	s, err := newJSONLSource(CacheFilename(ent))
	if err != nil {
		return nil, fmt.Errorf("The node %v doesn't have cache.\nPlease run all without cache again.", ent.NodeName)
	}
	if _, err := tb.Topology().AddSource(ent.NodeName, s, &core.SourceConfig{
		PausedOnStartup: true,
	}); err != nil {
		return nil, fmt.Errorf("canoot add a source for a cache file: %v", err)
	}
	return []string{ent.NodeName}, nil
}

func addStmt(tb *bql.TopologyBuilder, stmt *Statement) (core.Node, error) {
	n, err := tb.AddStmt(stmt.Stmt)
	if err != nil {
		return nil, fmt.Errorf(`cannot add a statement "%v": %v`, stmt, err)
	}
	return n, nil
}

// setUpState tries to load a saved state from the latest cache.
func setUpState(tb *bql.TopologyBuilder, latestEnt *CacheEntry, stmt *Statement) error {
	if latestEnt == nil { // the first entry
		_, err := addStmt(tb, stmt)
		return err
	}

	name := string(stmt.Stmt.(parser.CreateStateStmt).Name)
	typeName := string(stmt.Stmt.(parser.CreateStateStmt).Type)
	for _, n := range latestEnt.States {
		if name != n {
			continue
		}

		c, err := tb.UDSCreators.Lookup(typeName)
		if err != nil {
			return fmt.Errorf("cannot find the creator of the UDS type '%v': %v", typeName, err)
		}
		l, ok := c.(udf.UDSLoader)
		if !ok {
			break
		}

		err = func() error {
			fn := StateCacheFilename(latestEnt, name)
			f, err := os.Open(fn)
			if err != nil {
				return fmt.Errorf("cannot open the cache file of the UDS '%v': %v", err)
			}
			defer f.Close()
			r := bufio.NewReader(f)
			s, err := l.LoadState(tb.Topology().Context(), r, data.Map{}) // TODO: support loading parameters
			if err != nil {
				return fmt.Errorf("cannot load the saved UDS '%v' from the cache: %v", name, err)
			}
			return tb.Topology().Context().SharedStates.Add(name, typeName, s)
		}()
		return err
	}
	_, err := addStmt(tb, stmt)
	return err
}

func executeTarget(t core.Topology, ent *CacheEntry, stmt *Statement, resumeAll func() error) error {
	switch stmt.Stmt.(type) {
	case parser.CreateSourceStmt, parser.CreateStreamAsSelectStmt, parser.CreateStreamAsSelectUnionStmt:
		sink, err := newJSONLSink(CacheFilename(ent))
		if err != nil {
			return fmt.Errorf("cannot create a cache file %v: %v", CacheFilename(ent), err)
		}
		sn, err := t.AddSink("sensorbee_exp_jsonl_sink", sink, nil)
		if err != nil {
			return fmt.Errorf("cannot add a sink for a cache file: %v", err)
		}
		if err := sn.Input(ent.NodeName, nil); err != nil {
			return fmt.Errorf("cannot add an input from %v to a sink for a cache file: %v", ent.NodeName, err)
		}
		sn.StopOnDisconnect()
		if err := resumeAll(); err != nil {
			return err
		}
		sn.State().Wait(core.TSStopped)

	case parser.InsertIntoFromStmt, parser.InsertIntoSelectStmt:
		if err := resumeAll(); err != nil {
			return err
		}
		s, err := t.Sink(ent.NodeName)
		if err != nil {
			return fmt.Errorf("cannot find the sink %v: %v", ent.NodeName, err)
		}
		s.StopOnDisconnect()
		s.State().Wait(core.TSStopped)
	}
	return nil
}

func saveStates(t core.Topology, ent *CacheEntry) error {
	states, err := t.Context().SharedStates.List()
	if err != nil {
		return fmt.Errorf("cannot list created UDSs: %v", err)
	}
	for name, state := range states {
		s, ok := state.(core.SavableSharedState)
		if !ok {
			continue
		}

		err := func() error {
			fn := StateCacheFilename(ent, name)
			f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("cannot create a cache file for a UDS '%v': %v", fn, err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			if err := s.Save(t.Context(), w, data.Map{}); err != nil {
				return err
			}
			return w.Flush()
		}()
		if err != nil {
			return err
		}
		ent.States = append(ent.States, name)
	}
	return nil
}

func hashStatement(prev string, stmt *Statement) string {
	// TODO: compute hash taking dependencies into account so that cache can
	// tolerate reordering of statements. Be careful of some statements like
	// UPDATE which may have hidden effect on later execution.
	h := sha256.New()
	io.WriteString(h, prev)
	io.WriteString(h, stmt.String())
	return fmt.Sprintf("%x", h.Sum(nil))
}
