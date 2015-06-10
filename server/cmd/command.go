package cmd

// Command is the interface to management SensorBee.
type Command interface {
	// Init the command. Returns error if the command initialize is failed.
	Init() error
	// Name returns the names of the command. Users can execute the command
	// function by inputting these names.
	Name() []string
	// Execute the command function. Returns error if the function fail to
	// be executing. If input command is on the way, return false, and user
	// can add next command.
	Execute(input string) (bool, error)
}
