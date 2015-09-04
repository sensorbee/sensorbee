package shell

import (
	"flag"
	"fmt"
	"github.com/peterh/liner"
	"io"
	"os"
	"path"
	"pfi/sensorbee/sensorbee/client"
	"strings"
)

// App is the application server of SensorBee.
type App struct {
	historyFn            string
	executeExternalFiles *string
	requester            *client.Requester
	commandMap           map[string]Command
}

var tempDir = os.TempDir()

// SetUpCommands set up application. Commands are initialized with it.
func SetUpCommands(commands []Command) App {
	app := App{
		historyFn:            path.Join(tempDir, ".sensorbee_liner_history"),
		executeExternalFiles: flag.String("file", "", "execute BQL commands from external files"),
		commandMap:           map[string]Command{},
	}
	if len(commands) == 0 {
		return app
	}

	for _, cmd := range commands {
		cmd.Init()
		for _, v := range cmd.Name() {
			app.commandMap[v] = cmd
		}
	}
	return app
}

func (a *App) prompt(line *liner.State) {
	for {
		input, err := line.Prompt(promptLineStart)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "error reading line: %v", err)
			}
			return
		}

		if input != "" {
			line.AppendHistory(input)

			if strings.ToLower(input) == "exit" {
				fmt.Fprintln(os.Stdout, "SensorBee is closed")
				return
			}

			in := strings.ToLower(strings.Split(input, " ")[0])
			if cmd, ok := a.commandMap[in]; ok {
				a.processCommand(line, cmd, input)
			} else {
				fmt.Fprintf(os.Stderr, "undefined command: %v\n", in)
			}
		}
	}
}

func (a *App) processCommand(line *liner.State, cmd Command, input string) {
	for {
		if input != "" {
			line.AppendHistory(input)
			status, err := cmd.Input(input)
			switch status {
			case invalidCMD:
				fmt.Fprintf(os.Stderr, "input command is invalid: %v\n", err)
				return
			case preparedCMD:
				cmd.Eval(a.requester)
				return
			}
		}

		var err error
		input, err = line.Prompt(promptLineContinue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

// Run begins SensorBee command line tool to management SensorBee
// and execute BQL/UDF/UDTF statements.
func (a *App) Run(requester *client.Requester) {
	// set local value. For example, when stated with "-file=hoge.bql",
	// `executeExternalFile` will set "hoge.bql"
	flag.Parse()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(func(line string) (c []string) {
		// set auto complete command, if need
		return
	})

	if f, err := os.Open(a.historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	fmt.Fprintln(os.Stdout, appRunMsg)
	a.requester = requester
	a.prompt(line)

	if f, err := os.Create(a.historyFn); err != nil {
		fmt.Fprintf(os.Stderr, "error writing history file: %v", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

const (
	appRunMsg          = "SensorBee client tool is started!"
	promptLineStart    = ">>> "
	promptLineContinue = "... "
)
