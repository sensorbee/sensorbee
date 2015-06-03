package cmd

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"strings"
	"sync"
)

// state is manage line history.
// ref: https://github.com/peterh/liner
type state struct {
	r            *bufio.Reader
	history      []string
	historyMutex sync.RWMutex
}

func newState() *state {
	var s state
	s.r = bufio.NewReader(os.Stdin)
	return &s
}

func (s *state) promptUnsupported(p string) (string, error) {
	fmt.Fprint(os.Stdout, p)
	lineBuf, _, err := s.r.ReadLine()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(lineBuf)), nil
}

func (s *state) prompt(p string) (string, error) {
	return s.promptUnsupported(p)
}

func (s *state) close() error {
	return nil
}

func (s *state) appendHistory(item string) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()

	if len(s.history) > 0 {
		if item == s.history[len(s.history)-1] {
			return // TODO why?
		}
	}
	s.history = append(s.history, item)
	if len(s.history) > 1000 { // HistoryLimit is 1000
		s.history = s.history[1:]
	}
}

// liner is input string from os.stdin. Input strings are record to
// buffer, and will flush when one BQL statement is completed.
type liner struct {
	*state
	buffer string
}

func newLiner() *liner {
	s := newState()
	return &liner{
		state: s,
	}
}

func (li *liner) prefixInput() string {
	if li.buffer != "" {
		return promptLineContinue
	}
	return promptLineStart
}

func (li *liner) input() error {
	line, err := li.state.prompt(li.prefixInput())
	if err == nil {
		if li.buffer != "" {
			li.buffer = li.buffer + "\n" + line
		} else {
			li.buffer = line
		}
	}
	return err
}

func (li *liner) flush() {
	li.state.appendHistory(li.buffer)
	li.buffer = ""
}

func (li *liner) reinput() {
	lines := strings.Split(li.buffer, "\n")
	if strings.HasSuffix(lines[len(lines)-1], promptStatementEnd) {
		return
	}
	li.input()
	li.reinput()
}

func (li *liner) get() string {
	s := strings.Replace(li.buffer, "\n", " ", -1)
	s = s[:len(s)-1]
	return s
}

// Controller is a session control object to manage commands from external
// files or os.stdin
type controller struct {
	tb *bql.TopologyBuilder
	tp core.StaticTopology
}

func newController() (*controller, error) {
	c := &controller{
		tb: bql.NewTopologyBuilder(),
	}

	return c, nil
}

func (c *controller) readFiles(files []string) {
	for _, file := range files {
		bqls, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "external file reading error: %s", err.Error())
			break
		}
		if err = c.readFile(bqls); err != nil {
			fmt.Fprintln(os.Stderr, "external file reading error: %s", err.Error())
		}
	}
}

func (c *controller) readFile(src []byte) error {
	// TODO read file and input to controller
	return nil
}

func (c *controller) eval(in string) error {
	fmt.Printf("BQL: %s\n", in) // debug BQL
	err := c.tb.BQL(in)
	return err
}

func (c *controller) start() error {
	if c.tp == nil {
		tp, err := c.tb.Build()
		if err != nil {
			return err
		}
		c.tp = tp
	}
	// TODO create context
	logManager := core.NewConsolePrintLogger()
	conf := core.Configuration{
		TupleTraceEnabled: 1,
	}
	ctx := core.Context{
		Logger: logManager,
		Config: conf,
	}
	err := c.tp.Run(&ctx) // TODO use goroutine?
	fmt.Println("!!!!! topology is run and end !!!!!")
	return err
}

var executeExternalFiles = flag.String("file", "", "execute BQL commands from external files")

// Command starts SensorBee command line tool to execute BQL statements.
// This mode could stop typing "exit" or "fin" or "stop".
func CommandStart() {
	flag.Parse()

	c, err := newController()
	if err != nil {
		panic(err) // TODO panic is OK?
	}
	fmt.Fprintf(os.Stdout, "SensorBee %s\n", "0.1") // TODO get SensorBee version

	// external files
	if *executeExternalFiles != "" {
		externalFiles := strings.Split(*executeExternalFiles, ",")
		c.readFiles(externalFiles)
	}

	liner := newLiner()
	defer liner.close()
	for {
		err := liner.input()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal error has occurred: %s\n", err.Error())
			break
		}

		if liner.buffer == "" {
			continue
		}
		if _, ok := promptExitWords[liner.buffer]; ok {
			fmt.Fprintln(os.Stdout, "SensorBee command line tool was stopped")
			break
		}
		if _, ok := promptSensorBeeStartWords[liner.buffer]; ok {
			err = c.start()
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			}
			liner.flush()
			continue
		}

		liner.reinput()
		err = c.eval(liner.get())
		if err != nil {
			fmt.Println(err.Error())
		}
		liner.flush()
	}
}

const (
	promptLineStart    = ">>> "
	promptLineContinue = "... "
	promptStatementEnd = ";"
)

var promptExitWords = map[string]struct{}{
	"exit": struct{}{},
	"fin":  struct{}{},
	"stop": struct{}{},
}

var promptSensorBeeStartWords = map[string]struct{}{
	"run":   struct{}{},
	"start": struct{}{},
}
