// Copyright 2019 Tyler Yahn. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"io"
	"os"

	severity "github.com/MrAlias/go-logger/severity"
)

// std is the default logger.
var std = New(severity.Info, 0, os.Stderr)

// SetOutput sets the output destination for the default logger.
func SetOutput(w io.Writer) { std.SetOutput(w) }

// Flags returns the output flags for the default logger.
func Flags() int { return std.Flags() }

// SetFlags sets the output flags for the default logger.
func SetFlags(flags int) { std.SetFlags(flags) }

// Severity returns the minimum severity level the default logger logs at.
func Severity() severity.Level { return std.Severity() }

// SetSeverity sets the minimum severity level the default logger logs at.
func SetSeverity(level severity.Level) { std.SetSeverity(level) }

// Debug logs with `debug` severity to the default logger. Arguments are
// handled in the same manner as fmt.Print, but a new line is appended to
// the end if one not specified.
func Debug(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.debug.Output(2, fmt.Sprint(v...))
}

// Debugf logs with `debug` severity to the default logger. Arguments are
// handled in the same manner as fmt.Printf, but a new line is appended to
// the end if one not specified.
func Debugf(format string, v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.debug.Output(2, fmt.Sprintf(format, v...))
}

// Debugln logs with `debug` severity to the default logger. Arguments are
// handled in the same manner as fmt.Println.
func Debugln(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.debug.Output(2, fmt.Sprintln(v...))
}

// Info is like Debug, but logs with `info` severity.
func Info(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.info.Output(2, fmt.Sprint(v...))
}

// Infof is like Debugf, but logs with `info` severity.
func Infof(format string, v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.info.Output(2, fmt.Sprintf(format, v...))
}

// Infoln is like Debugln, but logs with `info` severity.
func Infoln(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.info.Output(2, fmt.Sprintln(v...))
}

// Error is like Debug, but logs with `error` severity.
func Error(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.error.Output(2, fmt.Sprint(v...))
}

// Errorf is like Debugf, but logs at an `error` severity.
func Errorf(format string, v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.error.Output(2, fmt.Sprintf(format, v...))
}

// Errorln is like Debugln, but logs with `error` severity.
func Errorln(v ...interface{}) {
	std.RLock()
	defer std.RUnlock()
	std.error.Output(2, fmt.Sprintln(v...))
}
