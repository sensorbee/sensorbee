package shell

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"io/ioutil"
	"strings"
)

// NewFileLoadCommands returns command list to load BQL file.
func NewFileLoadCommands() []Command {
	return []Command{
		&fileLoadCmd{},
	}
}

type fileLoadCmd struct {
	queries  string
	filePath string
}

func (f *fileLoadCmd) Init() error {
	return nil
}

func (f *fileLoadCmd) Name() []string {
	return []string{"file"}
}

func (f *fileLoadCmd) Input(input string) (cmdInputStatusType, error) {
	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return invalidCMD, fmt.Errorf("cannot read file name")
	}
	f.filePath = inputs[1]

	// file load
	queries, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return invalidCMD, fmt.Errorf("cannot read queries from file: %v", f.filePath)
	}
	f.queries = string(queries)
	return preparedCMD, nil
}

func (f *fileLoadCmd) Eval(requester *client.Requester) {
	fmt.Println("Eval() not implemented for fileLoadCmd, use App instead")
}
