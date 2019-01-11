// Copyright 2019 Tyler Yahn. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

var levelString = map[Level]string{
	Debug: "debug",
	Info:  "info",
	Error: "error",
}

// String returns severity level string representation.
func (l Level) String() string {
	return levelString[l]
}

// Prefix returns severity level logging prefix.
func (l Level) Prefix() string {
	return fmt.Sprintf("%-5s: ", strings.ToUpper(l.String()))
}
