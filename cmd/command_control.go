package cmd

import (
	"flag"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"strings"
)

var (
	historyFn            = "/tmp/.liner_history"
	executeExternalFiles = flag.String("file", "", "execute BQL commands from external files")
	commandMap           = map[string](func(string) (bool, error)){}
)

func setUpCommands() {
	commands := []Command{
		&bqlCmd{}, &bqlExecuter{}, &bqlStop{},
	}

	for _, cmd := range commands {
		cmd.Init()
		for _, v := range cmd.Name() {
			commandMap[v] = cmd.Execute
		}
	}
}

func prompt(line *liner.State) {
	input, err := line.Prompt(promptLineStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading line: %v", err)
		return
	}

	if input != "" {
		line.AppendHistory(input)

		if input == "exit" {
			fmt.Fprintln(os.Stdout, "SensorBee is closed")
			return
		}

		in := strings.ToLower(strings.Split(input, " ")[0])
		if cmd, ok := commandMap[in]; ok {
			isComplete, err := cmd(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if !isComplete {
				continuous(line, cmd)
			}
		} else {
			fmt.Fprintf(os.Stdout, "not found the command: %v\n", in)
		}
	}
	prompt(line)
}
func continuous(line *liner.State, cmd func(string) (bool, error)) {
	input, err := line.Prompt(promptLineContinue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if input != "" {
		line.AppendHistory(input)
		isComplete, err := cmd(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if isComplete {
			return
		}
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

	if f, err := os.Open(historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	fmt.Fprintln(os.Stdout, "SensorBee is started...")
	setUpCommands()
	prompt(line)

	if f, err := os.Create(historyFn); err != nil {
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
