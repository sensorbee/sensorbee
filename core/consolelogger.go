package core

import (
	"fmt"
	"log"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type ConsoleLogManager struct{}

func (l *ConsoleLogManager) Log(level LogLevel, msg string, a ...interface{}) {
	log.Printf("[%-7v] %v", level.String(), fmt.Sprintf(msg, a...))
}

func (l *ConsoleLogManager) DroppedTuple(t *tuple.Tuple, msg string, a ...interface{}) {
	log.Printf("[DROPPED] Tuple Batch ID is %v, %v", t.BatchID, fmt.Sprintf(msg, a...))
}
