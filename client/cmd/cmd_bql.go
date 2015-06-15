package cmd

import (
	"fmt"
	"strings"
)

const (
	bqlHeader = "/bql"
)

type bqlCommandState struct {
	tenant string
}

var (
	bqlCmdState = bqlCommandState{""}
)

// NewBQLCommands return command list to execute BQL statement.
func NewBQLCommands() []Command {
	return []Command{
		&tenantsCmd{},
		&changeTenantCmd{},
		&tenantCmd{},
		&bqlCmd{},
		&bqlExecuter{},
		&bqlStop{},
	}
}

type tenantsCmd struct{}

func (t *tenantsCmd) Init() error {
	return nil
}

func (t *tenantsCmd) Name() []string {
	return []string{"tenants"}
}

func (t *tenantsCmd) Eval(input string) (requestType, string, cmdInputStatusType) {
	uri := bqlHeader
	return getRequst, uri, preparedCMD
}

type changeTenantCmd struct{}

func (ct *changeTenantCmd) Init() error {
	return nil
}

func (ct *changeTenantCmd) Name() []string {
	return []string{"use"}
}

func (ct *changeTenantCmd) Eval(input string) (requestType, string, cmdInputStatusType) {
	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return otherRequest, "", invalidCMD
	}

	name := inputs[1]
	bqlCmdState.tenant = name
	uri := bqlHeader + "/" + name
	return postRequest, uri, preparedCMD
}

type tenantCmd struct{}

func (t *tenantCmd) Init() error {
	return nil
}

func (t *tenantCmd) Name() []string {
	return []string{"tenant"}
}

func (t *tenantCmd) Eval(input string) (requestType, string, cmdInputStatusType) {
	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return otherRequest, "", invalidCMD
	}

	name := inputs[1]
	uri := bqlHeader + "/" + name
	return getRequst, uri, preparedCMD
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
	return []string{"select", "create", "insert"}
}

// Eval resolves input command to BQL statement
func (b *bqlCmd) Eval(input string) (requestType, string, cmdInputStatusType) {
	if b.buffer == "" {
		b.buffer = input
	} else {
		b.buffer += "\n" + input
	}
	if !strings.HasSuffix(input, ";") {
		return otherRequest, "", continuousCMD
	}

	// flush buffer and get complete statement
	stmt := strings.Replace(b.buffer, "\n", " ", -1)
	stmt = stmt[:len(stmt)-1]
	b.buffer = ""

	fmt.Printf("BQL: %s\n", stmt) // for debug, delete later

	uri := bqlHeader + "/" + bqlCmdState.tenant + "/" + stmt
	return postRequest, uri, preparedCMD
}

type bqlExecuter struct {
}

// Init (nothing to do)
func (be *bqlExecuter) Init() error {
	return nil
}

// Name returns topology start words.
func (be *bqlExecuter) Name() []string {
	return []string{"run", "start"}
}

// Eval topology run.
// TODO execute can only deal with static topology
// we have to deal with dynamic topology
func (be *bqlExecuter) Eval(input string) (requestType, string, cmdInputStatusType) {
	uri := bqlHeader + "/" + bqlCmdState.tenant + "/run"
	return postRequest, uri, preparedCMD
}

type bqlStop struct {
}

// Init (nothing to do)
func (bs *bqlStop) Init() error {
	return nil
}

// Name returns topology stop word.
func (bs *bqlStop) Name() []string {
	return []string{"stop"}
}

// Execute to stop topology, will be stopped all components.
// This means not to stop SensorBee.
func (bs *bqlStop) Eval(input string) (requestType, string, cmdInputStatusType) {
	uri := bqlHeader + "/" + bqlCmdState.tenant + "/stop"
	return postRequest, uri, preparedCMD
}
