package shell

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
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
	currentTopology currentTopologyState
)

// NewBQLCommands return command list to execute BQL statement.
func NewTopologiesCommands() []Command {
	return []Command{
		&changeTopologyCmd{},
		&bqlCmd{},
	}
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

func (ct *changeTopologyCmd) Eval(requester *client.Requester) {
	currentTopology.name = ct.name
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
	return []string{"select", "create", "insert", "resume", "update"}
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
func (b *bqlCmd) Eval(requester *client.Requester) {
	// flush buffer and get complete statement
	queries := b.buffer
	b.buffer = ""

	fmt.Printf("BQL: %s\n", queries) // for debug, delete later
	sendBQLQueries(requester, queries)
}

func sendBQLQueries(requester *client.Requester, queries string) {
	uri := topologiesHeader + "/" + currentTopology.name + "/queries"
	res, err := requester.Do(client.Post, uri, map[string]interface{}{
		"queries": queries,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
		return
	}
	defer res.Close()

	if res.IsError() {
		// TODO: provide error reporting utility
		errRes, err := res.Error()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// TODO: enhance error message
		fmt.Fprintf(os.Stderr, "request failed: %v: %v: %v\n", errRes.Code, errRes.Message, errRes.Meta)
		return
	}

	if res.IsStream() {
		showStreamResponses(res)
		return
	}
	// TODO: there isn't much information to show right now. Improve the server's response.
}

func showStreamResponses(res *client.Response) {
	ch, err := res.ReadStreamJSON()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	for {
		select {
		case js, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(js)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot marshal a JSON: %v\n", err)
				return
			}
			fmt.Printf("%s\n", data)

		case <-sig:
			return // The response is closed by the caller
		}
	}
}
