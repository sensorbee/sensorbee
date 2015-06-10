package cmd

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"strings"
)

// NewBQLCommands return command list to execute BQL statement.
func NewBQLCommands() []Command {
	return []Command{&bqlCmd{}, &bqlExecuter{}, &bqlStop{}}
}

type bqlCommandState struct {
	tb  *bql.TopologyBuilder
	tp  core.StaticTopology
	ctx core.Context
}

var bqlCmdState = bqlCommandState{}

type bqlCmd struct {
	buffer string
}

// Init BQL state.
func (b *bqlCmd) Init() error {
	bqlCmdState.tb = bql.NewTopologyBuilder()
	return nil
}

// Name returns BQL start words.
func (b *bqlCmd) Name() []string {
	return []string{"select", "create", "insert"}
}

// Execute topology builder building, and register BQL statements.
func (b *bqlCmd) Execute(input string) (bool, error) {
	if b.buffer == "" {
		b.buffer = input
	} else {
		b.buffer += "\n" + input
	}
	if !strings.HasSuffix(input, ";") {
		return false, nil
	}

	// flush buffer and get complete statement
	stmt := strings.Replace(b.buffer, "\n", " ", -1)
	stmt = stmt[:len(stmt)-1]
	b.buffer = ""

	fmt.Printf("BQL: %s\n", stmt) // for debug, delete later
	err := bqlCmdState.tb.BQL(stmt)
	return true, err
}

type bqlExecuter struct {
}

// Init (nothing to do)
func (be *bqlExecuter) Init() error {
	return nil
}

// Name returns topology start words.
func (be *bqlExecuter) Name() []string {
	return []string{"run", "start"}
}

// Execute topology run.
// TODO execute can only deal with static topology
// we have to deal with dynamic topology
func (be *bqlExecuter) Execute(input string) (bool, error) {
	if bqlCmdState.tp == nil {
		tp, err := bqlCmdState.tb.Build()
		if err != nil {
			return true, err
		}
		bqlCmdState.tp = tp

		// TODO create context
		logManaget := core.NewConsolePrintLogger()
		conf := core.Configuration{
			TupleTraceEnabled: 1,
		}
		ctx := core.Context{
			Logger: logManaget,
			Config: conf,
		}
		bqlCmdState.ctx = ctx
	}

	err := bqlCmdState.tp.Run(&(bqlCmdState.ctx))  // TODO use goroutine?
	fmt.Println("!!! topology is run and end !!!") // for debug, delete later
	return true, err
}

type bqlStop struct {
}

// Init (nothing to do)
func (bs *bqlStop) Init() error {
	return nil
}

// Name returns topology stop word.
func (bs *bqlStop) Name() []string {
	return []string{"stop"}
}

// Execute to stop topology, will be stopped all components.
// This means not to stop SensorBee.
func (bs *bqlStop) Execute(input string) (bool, error) {
	if bqlCmdState.tp == nil {
		return true, fmt.Errorf("not set up topology: %v", "")
	}
	err := bqlCmdState.tp.Stop(&(bqlCmdState.ctx))
	return true, err
}
