// Copyright 2019 Tyler Yahn. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	severity "github.com/MrAlias/go-logger/severity"
)

// Logger - Severity based log writer.
type Logger struct {
	sync.RWMutex

	out      io.Writer
	severity severity.Level
	flags    int

	debug *log.Logger
	info  *log.Logger
	error *log.Logger
}

// New returns a new Logger set to log with severity level and flags to w.
func New(level severity.Level, flags int, w io.Writer) *Logger {
	l := &Logger{}
	l.init(level, flags, w)
	return l
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.init(l.severity, l.flags, w)
}

// Flags returns the output flags for the logger.
func (l *Logger) Flags() int {
	l.RLock()
	defer l.RUnlock()
	return l.flags
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flags int) {
	l.init(l.severity, flags, l.out)
}

// Severity returns the minimum severity level logged.
func (l *Logger) Severity() severity.Level {
	l.RLock()
	defer l.RUnlock()
	return l.severity
}

// SetSeverity sets minimum severity level logged.
func (l *Logger) SetSeverity(level severity.Level) {
	l.init(level, l.flags, l.out)
}

// init initializes the logging handlers of l in a concurrency-safe manner.
func (l *Logger) init(level severity.Level, flags int, w io.Writer) {
	setLoggers := func(dW, iW, eW io.Writer) (*log.Logger, *log.Logger, *log.Logger) {
		dP := severity.Debug.Prefix()
		iP := severity.Info.Prefix()
		eP := severity.Error.Prefix()
		return log.New(dW, dP, flags), log.New(iW, iP, flags), log.New(eW, eP, flags)
	}

	l.Lock()
	defer l.Unlock()

	switch level {
	case severity.Debug:
		l.debug, l.info, l.error = setLoggers(w, w, w)
	case severity.Info:
		l.debug, l.info, l.error = setLoggers(ioutil.Discard, w, w)
	case severity.Error:
		l.debug, l.info, l.error = setLoggers(ioutil.Discard, ioutil.Discard, w)
	default:
		panic(fmt.Sprintf("invalid logger severity: %s", level))
	}

	l.out = w
	l.flags = flags
	l.severity = level

	return
}

// Debug logs with `debug` severity. Arguments are handled in the
// same manner as fmt.Print, but a new line is appended to the end
// if one not specified.
func (l *Logger) Debug(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.debug.Output(2, fmt.Sprint(v...))
}

// Debugf logs with `debug` severity. Arguments are handled in the same
// manner as fmt.Printf, but a new line is appended to the end if one not
// specified.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.debug.Output(2, fmt.Sprintf(format, v...))
}

// Debugln logs with `debug` severity. Arguments are handled in the same
// manner as fmt.Println.
func (l *Logger) Debugln(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.debug.Output(2, fmt.Sprintln(v...))
}

// Info is like Debug, but logs with `info` severity.
func (l *Logger) Info(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.info.Output(2, fmt.Sprint(v...))
}

// Infof is like Debugf, but logs with `info` severity.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.info.Output(2, fmt.Sprintf(format, v...))
}

// Infoln is like Debugln, but logs with `info` severity.
func (l *Logger) Infoln(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.info.Output(2, fmt.Sprintln(v...))
}

// Error is like Debug, but logs with `error` severity.
func (l *Logger) Error(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.error.Output(2, fmt.Sprint(v...))
}

// Errorf is like Debugf, but logs at an `error` severity.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.error.Output(2, fmt.Sprintf(format, v...))
}

// Errorln is like Debugln, but logs with `error` severity.
func (l *Logger) Errorln(v ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	l.error.Output(2, fmt.Sprintln(v...))
}
