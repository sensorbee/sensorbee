package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"pfi/sensorbee/sensorbee/tuple"
)

// NewConsolePrintLogger creates a LogManager that will write
// all passed in messages to stderr (just as the normal log
// modules does).
func NewConsolePrintLogger() LogManager {
	// same as log.std logger
	return NewPrintLogger(os.Stderr)
}

// NewPrintLogger creates a LogManager that will write all passed
// in messages to the given writer.
func NewPrintLogger(w io.Writer) LogManager {
	logger := log.New(w, "", log.LstdFlags)
	return &simplePrintLogger{logger}
}

type simplePrintLogger struct {
	logger *log.Logger
}

func (l *simplePrintLogger) Log(level LogLevel, msg string, a ...interface{}) {
	l.logger.Printf("[%-7v] %v", level.String(), fmt.Sprintf(msg, a...))
}

func (l *simplePrintLogger) DroppedTuple(t *tuple.Tuple, msg string, a ...interface{}) {
	l.logger.Printf("[DROPPED] Tuple Batch ID is %v, %v", t.BatchID, fmt.Sprintf(msg, a...))
}
