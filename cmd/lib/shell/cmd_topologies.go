package shell

import (
	"fmt"
	"pfi/sensorbee/sensorbee/client"
	"strings"
)

const (
	topologiesHeader = "/topologies"
)

type currentTopologyState struct {
	name string
}

var (
	currentTopology = currentTopologyState{""}
)

// NewBQLCommands return command list to execute BQL statement.
func NewTopologiesCommands() []Command {
	return []Command{
		&topologiesCmd{},
		&changeTopologyCmd{},
		&topologyCmd{},
		&topologyStopCmd{},
		&bqlCmd{},
	}
}

type topologiesCmd struct {
	uri string
}

func (t *topologiesCmd) Init() error {
	return nil
}

func (t *topologiesCmd) Name() []string {
	return []string{"topologies"}
}

func (t *topologiesCmd) Input(input string) (cmdInputStatusType, error) {
	t.uri = topologiesHeader
	return preparedCMD, nil
}

func (t *topologiesCmd) Eval() (client.Method, string, interface{}) {
	return client.Get, t.uri, nil
}

type changeTopologyCmd struct {
	name string
}

func (ct *changeTopologyCmd) Init() error {
	return nil
}

func (ct *changeTopologyCmd) Name() []string {
	return []string{"use"}
}

func (ct *changeTopologyCmd) Input(input string) (cmdInputStatusType, error) {
	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return invalidCMD, fmt.Errorf("cannot support empty named topology")
	}

	ct.name = inputs[1]
	return preparedCMD, nil
}

func (ct *changeTopologyCmd) Eval() (client.Method, string, interface{}) {
	currentTopology.name = ct.name
	return client.OtherMethod, "", nil
}

type topologyCmd struct {
	uri string
}

func (t *topologyCmd) Init() error {
	return nil
}

func (t *topologyCmd) Name() []string {
	return []string{"info"}
}

func (t *topologyCmd) Input(input string) (cmdInputStatusType, error) {
	inputs := strings.Split(input, " ")
	var name string
	if len(inputs) != 2 {
		if currentTopology.name == "" {
			return invalidCMD, fmt.Errorf("target topology is empty")
		}
		name = currentTopology.name
	} else {
		name = inputs[1]
	}

	t.uri = topologiesHeader + "/" + name
	return preparedCMD, nil
}

func (t *topologyCmd) Eval() (client.Method, string, interface{}) {
	return client.Get, t.uri, nil
}

type topologyStopCmd struct {
	uri string
}

// Init (nothing to do)
func (be *topologyStopCmd) Init() error {
	return nil
}

// Name returns topology stop words.
func (be *topologyStopCmd) Name() []string {
	return []string{"stop"}
}

func (be *topologyStopCmd) Input(input string) (cmdInputStatusType, error) {
	return preparedCMD, nil
}

// Eval operates topology stop.
func (be *topologyStopCmd) Eval() (client.Method, string, interface{}) {
	uri := topologiesHeader + "/" + currentTopology.name
	m := map[string]interface{}{}
	m["state"] = "stop"
	return client.Put, uri, &m
}

type bqlCmd struct {
	buffer string
}

// Init BQL state.
func (b *bqlCmd) Init() error {
	return nil
}

// Name returns BQL start words.
func (b *bqlCmd) Name() []string {
	return []string{"select", "create", "insert", "resume"}
}

func (b *bqlCmd) Input(input string) (cmdInputStatusType, error) {
	if b.buffer == "" {
		b.buffer = input
	} else {
		b.buffer += "\n" + input
	}
	if !strings.HasSuffix(input, ";") {
		return continuousCMD, nil
	}

	return preparedCMD, nil
}

// Eval resolves input command to BQL statement
func (b *bqlCmd) Eval() (client.Method, string, interface{}) {
	// flush buffer and get complete statement
	queries := strings.Replace(b.buffer, "\n", " ", -1)
	queries = queries[:len(queries)-1]
	b.buffer = ""

	fmt.Printf("BQL: %s\n", queries) // for debug, delete later

	uri := topologiesHeader + "/" + currentTopology.name + "/queries"
	m := map[string]interface{}{}
	m["queries"] = queries
	return client.Post, uri, &m
}
