package logging

import (
	"fmt"
	"time"
)

type Logger struct {
	Prefix string
}

func Local(prefix string) Logger {
	return Logger{Prefix: prefix}
}

func (logger Logger) Info(format string, args ...any) {
	fmt.Printf(time.Now().Format(time.StampMilli)+"| "+logger.Prefix+": "+format+"\n", args...)
}

func (logger Logger) Child(prefix string) Logger {
	return Logger{Prefix: logger.Prefix + "." + prefix}
}
