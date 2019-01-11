package logger

import (
	"fmt"
	"strings"
)

type Level int

// Levels of severity
const (
	Debug Level = iota
	Info
	Error
)

var levelsString = map[Level]string{
	Debug: "debug",
	Info:  "info",
	Error: "error",
}

// String returns severity level string representation.
func (l Level) String() string {
	return levelsString[l]
}

// Prefix returns severity level logging prefix.
func (l Level) Prefix() string {
	return fmt.Sprintf("%-5s: ", strings.ToUpper(l.String()))
}
