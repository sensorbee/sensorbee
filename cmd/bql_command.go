package cmd

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"strings"
)

type bqlCmdState struct {
	tb  *bql.TopologyBuilder
	tp  core.StaticTopology
	ctx core.Context
}

// command is the interface to management SensorBee.
type command interface {
	// init the command. Returns error if the command fail to initialize.
	init(state *bqlCmdState) error
	// name returns the names of the command. Users can execute the command
	// function by inputting these names. The command names are not
	// distinguished upper case or lower case.
	name() []string
	// execute the command function. Returns error if the function fail to
	// be executing. state have share object needed to execute the command.
	// If input command is on the way, return false, and user can input
	// next command.
	execute(state *bqlCmdState, input string) (bool, error)
}

type bqlCmd struct {
	buffer string
}

func (b *bqlCmd) init(state *bqlCmdState) error {
	state.tb = bql.NewTopologyBuilder()
	return nil
}

func (b *bqlCmd) name() []string {
	return []string{"select", "create", "insert"}
}

func (b *bqlCmd) execute(state *bqlCmdState, input string) (bool, error) {
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
	err := state.tb.BQL(stmt)
	return true, err
}

type bqlExecuter struct {
}

func (be *bqlExecuter) init(state *bqlCmdState) error {
	return nil
}

func (be *bqlExecuter) name() []string {
	return []string{"run", "start"}
}

// TODO execute can only deal with static topology
// we have to deal with dynamic topology
func (be *bqlExecuter) execute(state *bqlCmdState, input string) (bool, error) {
	if state.tp == nil {
		tp, err := state.tb.Build()
		if err != nil {
			return true, err
		}
		state.tp = tp

		// TODO create context
		logManaget := core.NewConsolePrintLogger()
		conf := core.Configuration{
			TupleTraceEnabled: 1,
		}
		ctx := core.Context{
			Logger: logManaget,
			Config: conf,
		}
		state.ctx = ctx
	}

	err := state.tp.Run(&(state.ctx))              // TODO use goroutine?
	fmt.Println("!!! topology is run and end !!!") // for debug, delete later
	return true, err
}

type bqlStop struct {
}

func (bs *bqlStop) init(state *bqlCmdState) error {
	return nil
}

func (bs *bqlStop) name() []string {
	return []string{"stop"}
}

func (bs *bqlStop) execute(state *bqlCmdState, input string) (bool, error) {
	if state.tp == nil {
		return true, fmt.Errorf("not set up topology: %v", "")
	}
	err := state.tp.Stop(&(state.ctx))
	return true, err
}
