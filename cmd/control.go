package cmd

import (
	"flag"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"strings"
)

var (
	history_fn           = "/tmp/.liner_history"
	executeExternalFiles = flag.String("file", "", "execute BQL commands from external files")
	commandState         = bqlCmdState{}
	commandMap           = map[string](func(bqlCmdState, string) (bool, error)){}
)

func setUpCommands() {
	commands := []command{
		&bqlCmd{}, &bqlExecuter{}, &bqlStop{},
	}

	for _, cmd := range commands {
		cmd.init(commandState)
		for _, v := range cmd.name() {
			commandMap[v] = cmd.execute
		}
	}
}

func prompt(line *liner.State) {
	input, err := line.Prompt(promptLineStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading line: %v", err)
		return
	}

	fmt.Fprintf(os.Stdout, "got: %v", input) // delete later
	line.AppendHistory(input)

	if input == "exit" {
		fmt.Fprintln(os.Stdout, "SensorBee is closed")
		return
	}

	in := strings.ToLower(strings.Split(input, " ")[0])
	if cmd, ok := commandMap[in]; ok {
		isContinue, err := cmd(commandState, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			prompt(line)
		}
		if isContinue {
			continuous(line, cmd)
		}
	}
}
func continuous(line *liner.State, cmd func(bqlCmdState, string) (bool, error)) {
	input, err := line.Prompt(promptLineContinue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	line.AppendHistory(input)
	isContinue, err := cmd(commandState, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !isContinue {
		return
	}
	continuous(line, cmd)
}

// CommandStart begins SensorBee command line tool to management SensorBee
// and execute BQL/UDF/UDTF statements.
func CommandStart() {
	// set local value. For example, when stated with "-file=hoge.bql",
	// `executeExternalFile` will set "hoge.bql"
	flag.Parse()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(func(line string) (c []string) {
		// set auto complete command, if need
		return
	})

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	setUpCommands()
	prompt(line)

	if f, err := os.Create(history_fn); err != nil {
		fmt.Fprintf(os.Stderr, "error writing history file: %v", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

const (
	promptLineStart    = ">>> "
	promptLineContinue = "... "
)
