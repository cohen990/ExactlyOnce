package logging

import (
	"fmt"
	"time"
)

type LogRoot struct {
	// Prefix is what goes in front of the string when it logs.
	Prefix string
}

type Logger struct {
	LogRoot
}

func NewRoot(prefix string) LogRoot {
	return LogRoot{Prefix: prefix}
}

func (logger Logger) Info(format string, args ...any) {
	fmt.Printf(time.Now().Format(time.StampMilli)+" | "+logger.Prefix+": "+format+"\n", args...)
}

func (root LogRoot) Child(prefix string) Logger {
	return Logger{LogRoot{prefix}}
}
