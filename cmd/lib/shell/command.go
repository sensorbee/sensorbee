package shell

import (
	"pfi/sensorbee/sensorbee/client"
)

type cmdInputStatusType int

const (
	// prepared means input command is valid and converted to URL.
	preparedCMD cmdInputStatusType = iota
	// continuous means input command is not completed, and user can
	// input continued commands.
	continuousCMD
	// invalid means input command can not resolve on converting URL.
	invalidCMD
)

// Command is the interface to management SensorBee.
type Command interface {
	// Init the command. Returns error if the command initialize is failed.
	Init() error
	// Name returns the names of the command. Users can execute the command
	// function by inputting these names.
	Name() []string
	// Input commands to buffer. If commands are completed, returns that
	// cmdInputStatusType is preparedCMD, and commands are on the way,
	// returns that cmdInputStatusType is continuousCMD.
	// Returns error when the input commands are invalid.
	Input(input string) (cmdInputStatusType, error)
	// Eval resolve input command to convert URL and requestType.
	Eval() (client.Method, string, interface{})
}
