package shell

import (
	"bufio"
	"fmt"
	"github.com/peterh/liner"
	"io"
	"os"
	"path"
	"pfi/sensorbee/sensorbee/client"
	"runtime"
	"strings"
)

// App is the application server of SensorBee.
type App struct {
	historyFn  string
	requester  *client.Requester
	commandMap map[string]Command
}

// SetUpCommands set up application. Commands are initialized with it.
func SetUpCommands(commands []Command) App {
	// get home directory for the history file
	histDir := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		histDir = path.Join(os.Getenv("APPDATA"), "PFN", "SensorBee")
		if err := os.MkdirAll(histDir, os.ModeDir); err != nil {
			fmt.Printf("failed to create application directory '%s': %s\n", histDir, err)
		}
	}

	app := App{
		historyFn:  path.Join(histDir, ".sensorbee_liner_history"),
		commandMap: map[string]Command{},
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
	// create a function to read from the terminal
	// and output an appropriate prompt
	// (if `continued` is true, then the prompt
	// for a continued statement will be shown)
	getNextLine := func(continued bool) (string, error) {
		// create prompt
		promptStart := "(no topology)" + promptLineStart
		if currentTopology.name != "" {
			promptStart = fmt.Sprintf("(%s)%s", currentTopology.name,
				promptLineStart)
		}
		if continued {
			promptStart = ""
		}

		// get line from terminal
		input, err := line.Prompt(promptStart)
		if err == nil && input != "" {
			line.AppendHistory(input)
		}
		return input, err
	}

	// continue as long as there is input
	for a.readStartOfNextCommand(getNextLine, true) {
	}
}

func (a *App) readStartOfNextCommand(getNextLine func(bool) (string, error), topLevel bool) bool {
	input, err := getNextLine(false)
	// if there is no next line, stop
	if err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "error reading line: %v", err)
		}
		if topLevel {
			// there was an EOF control character, e.g., the
			// user pressed Ctrl+D. in order not to mess up the
			// terminal, write an additional newline character
			fmt.Println("")
		}
		return false
	}

	// if there is input, find the type of command that was input
	if input != "" {
		if strings.ToLower(input) == "exit" {
			if !topLevel {
				fmt.Fprintln(os.Stdout, "exit from file processing")
			}
			return false
		}

		if strings.HasPrefix(input, "--") { // BQL comment
			return true
		}

		// detect which type of command we have
		in := strings.ToLower(strings.Split(input, " ")[0])
		if cmd, ok := a.commandMap[in]; ok {
			// continue to read until the end of the statement
			a.readCompleteCommand(cmd, input, getNextLine)
		} else {
			fmt.Fprintf(os.Stderr, "undefined command: %v\n", in)
		}
	}
	return true
}

func (a *App) readCompleteCommand(cmd Command, input string, getNextLine func(bool) (string, error)) {
	for {
		if input != "" {
			// feed string to the command we detected and check
			// if the command accepts this as valid continuation
			status, err := cmd.Input(input)
			switch status {
			case invalidCMD:
				fmt.Fprintf(os.Stderr, "input command is invalid: %v\n", err)
				return
			case preparedCMD:
				// if the statement is complete, evaluate it
				if fileCmd, ok := cmd.(*fileLoadCmd); ok {
					r := strings.NewReader(fileCmd.queries)
					b := bufio.NewReader(r)
					getNextLineInFile := func(continued bool) (string, error) {
						s, err := b.ReadString('\n')
						return strings.TrimSuffix(s, "\n"), err
					}
					// continue as long as there is input
					path := fileCmd.filePath
					fmt.Printf("-- process %s --\n", path)
					for a.readStartOfNextCommand(getNextLineInFile, false) {
					}
					fmt.Printf("-- end %s --\n", path)
					return
				}
				cmd.Eval(a.requester)
				return
			}
		}

		var err error
		input, err = getNextLine(true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

// Run begins SensorBee command line tool to management SensorBee
// and execute BQL/UDF/UDTF statements.
func (a *App) Run(requester *client.Requester) {
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
	promptLineStart = "> "
)
