// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"

	severity "github.com/MrAlias/go-logger/severity"
)

// Sourced from the std log package:
// https://golang.org/src/log/log_test.go
const (
	rDate         = `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	rTime         = `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	rMicroseconds = `\.[0-9][0-9][0-9][0-9][0-9][0-9]`
	rLine         = `(92|96|100|110|114|118|123|127|131):` // must update if the calls to log below move
	rLongfile     = `.*/[A-Za-z0-9_\-]+\.go:` + rLine
	rShortfile    = `[A-Za-z0-9_\-]+\.go:` + rLine
)

func validate(t *testing.T, logFunc, expected, line string, notEmpty bool) {
	if notEmpty {
		line = line[0 : len(line)-1] // skip the newline on non-empty logs
	}
	if matched, err := regexp.MatchString(expected, line); err != nil {
		t.Fatalf("pattern %q failed to compile: %v", expected, err)
	} else if !matched {
		t.Errorf("%s log output should match %q is %q", logFunc, expected, line)
	}
}

func testLoggerSeverity(t *testing.T, l *Logger, flag int, pattern string, level severity.Level) {
	origFlags := l.Flags()
	origSeverity := l.Severity()

	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.SetFlags(flag)
	if f := l.Flags(); f != flag {
		t.Errorf("logger flags should match %v is %v", flag, f)
	}

	l.SetSeverity(level)
	if s := l.Severity(); s != level {
		t.Errorf("logger severity should match %q is %q", level.String(), s)
	}

	noPrefixPattern := "^%s" + pattern + "hello 23 world$"
	var expected string

	// No for loop like testStdLoggerSeverity so as to keep the call stack at 2.
	if severity.Debug >= level {
		expected = fmt.Sprintf(noPrefixPattern, severity.Debug.Prefix())
	} else {
		expected = "^$"
	}

	buf.Reset()
	l.Debug("hello 23 world")
	validate(t, "l.Debug", expected, buf.String(), severity.Debug >= level)

	buf.Reset()
	l.Debugf("hello %d world", 23)
	validate(t, "l.Debugf", expected, buf.String(), severity.Debug >= level)

	buf.Reset()
	l.Debugln("hello", 23, "world")
	validate(t, "l.Debugln", expected, buf.String(), severity.Debug >= level)

	if severity.Info >= level {
		expected = fmt.Sprintf(noPrefixPattern, severity.Info.Prefix())
	} else {
		expected = "^$"
	}

	buf.Reset()
	l.Info("hello 23 world")
	validate(t, "l.Info", expected, buf.String(), severity.Info >= level)

	buf.Reset()
	l.Infof("hello %d world", 23)
	validate(t, "l.Infof", expected, buf.String(), severity.Info >= level)

	buf.Reset()
	l.Infoln("hello", 23, "world")
	validate(t, "l.Infoln", expected, buf.String(), severity.Info >= level)

	expected = fmt.Sprintf(noPrefixPattern, severity.Error.Prefix())
	buf.Reset()
	l.Error("hello 23 world")
	validate(t, "l.Error", expected, buf.String(), true)

	buf.Reset()
	l.Errorf("hello %d world", 23)
	validate(t, "l.Errorf", expected, buf.String(), true)

	buf.Reset()
	l.Errorln("hello", 23, "world")
	validate(t, "l.Errorln", expected, buf.String(), true)

	SetSeverity(origSeverity)
	SetFlags(origFlags)
	SetOutput(ioutil.Discard)
}

func TestAllLoggerLogging(t *testing.T) {
	var tests = []struct {
		flag    int
		pattern string // regexp that log output must match; we add `^`, `prefix`, and `expected_text$` always
	}{
		{0, ""},
		{log.Ldate, rDate + " "},
		{log.Ltime, rTime + " "},
		{log.Ltime | log.Lmicroseconds, rTime + rMicroseconds + " "},
		{log.Lmicroseconds, rTime + rMicroseconds + " "}, // microsec implies time
		{log.Llongfile, rLongfile + " "},
		{log.Lshortfile, rShortfile + " "},
		{log.Llongfile | log.Lshortfile, rShortfile + " "}, // shortfile overrides longfile
		{log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile, rDate + " " + rTime + rMicroseconds + " " + rLongfile + " "},
		{log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile, rDate + " " + rTime + rMicroseconds + " " + rShortfile + " "},
	}

	l := New(severity.Info, 0, ioutil.Discard)
	for _, test := range tests {
		testLoggerSeverity(t, l, test.flag, test.pattern, severity.Debug)
		testLoggerSeverity(t, l, test.flag, test.pattern, severity.Info)
		testLoggerSeverity(t, l, test.flag, test.pattern, severity.Error)
	}
}
