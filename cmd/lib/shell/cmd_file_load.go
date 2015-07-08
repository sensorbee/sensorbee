package shell

import (
	"fmt"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/client"
	"strings"
)

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

func (f *fileLoadCmd) Eval() (client.Method, string, interface{}) {
	uri := topologiesHeader + "/" + currentTopology.name + "/queries"
	m := map[string]interface{}{}
	m["queries"] = f.queries
	return client.Post, uri, m
}
