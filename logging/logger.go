package logging

import (
	"fmt"
	"time"
)

type LogRoot struct {
	// Prefix is what goes in front of the string when it logs.
	Prefix string
	muted  bool
}

type Logger struct {
	LogRoot
}

func NewRoot(prefix string) LogRoot {
	return LogRoot{Prefix: prefix}
}

func (logRoot LogRoot) Mute() LogRoot {
	logRoot.muted = true
	return logRoot
}

func (logger Logger) Info(format string, args ...any) {
	if logger.muted {
		return
	}

	fmt.Printf(time.Now().Format(time.StampMilli)+" | "+logger.Prefix+": "+format+"\n", args...)
}

func (root LogRoot) Child(prefix string) Logger {
	return Logger{LogRoot{Prefix: prefix}}
}
