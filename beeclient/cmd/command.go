package cmd

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
	// Eval resolve input command to convert URL. Returns cmdInputStatus that
	// result of resolving the input command.
	Eval(input string) (requestType, string, cmdInputStatusType)
}
