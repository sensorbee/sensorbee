package shell

import (
	"fmt"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/client"
	"strings"
)

// NewFileLoadCommands returns command list to load BQL file.
func NewFileLoadCommands() []Command {
	return []Command{
		&fileLoadCmd{},
	}
}

type fileLoadCmd struct {
	queries string
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
	filePath := inputs[1]

	// file load
	queries, err := ioutil.ReadFile(filePath)
	if err != nil {
		return invalidCMD, fmt.Errorf("cannot read queries from file: %v", filePath)
	}
	f.queries = string(queries)
	return preparedCMD, nil
}

func (f *fileLoadCmd) Eval(requester *client.Requester) {
	fmt.Println("Eval() not implemented for fileLoadCmd, use App instead")
}
