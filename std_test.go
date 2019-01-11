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
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	severity "github.com/MrAlias/go-logger/severity"
)

// Sourced from the std log package:
// https://golang.org/src/log/log_test.go
const (
	rStdDate         = `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	rStdTime         = `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	rStdMicroseconds = `\.[0-9][0-9][0-9][0-9][0-9][0-9]`
	rStdLine         = `(102|106|110):` // must update if the calls to log below move
	rStdLongfile     = `.*/[A-Za-z0-9_\-]+\.go:` + rStdLine
	rStdShortfile    = `[A-Za-z0-9_\-]+\.go:` + rStdLine
)

func validateStd(t *testing.T, logFunc, expected, line string, notEmpty bool) {
	if notEmpty {
		line = line[0 : len(line)-1] // skip the newline on non-empty logs
	}
	if matched, err := regexp.MatchString(expected, line); err != nil {
		t.Fatalf("pattern %q failed to compile: %v", expected, err)
	} else if !matched {
		t.Errorf("%s log output should match %q is %q", logFunc, expected, line)
	}
}

func testStdLoggerSeverity(t *testing.T, flag int, pattern string, level severity.Level) {
	buf := new(bytes.Buffer)
	SetOutput(buf)

	SetFlags(flag)
	if f := Flags(); f != flag {
		t.Errorf("std logger flags should match %v is %v", flag, f)
	}

	SetSeverity(level)
	if s := Severity(); s != level {
		t.Errorf("std logger severity should match %q is %q", level.String(), s)
	}

	var severityTests = []struct {
		Severity    severity.Level
		PrintFunc   func(...interface{})
		PrintfFunc  func(string, ...interface{})
		PrintlnFunc func(...interface{})
	}{
		{severity.Debug, Debug, Debugf, Debugln},
		{severity.Info, Info, Infof, Infoln},
		{severity.Error, Error, Errorf, Errorln},
	}

	noPrefixPattern := "^%s" + pattern + "hello 23 world$"
	for _, test := range severityTests {
		var expected string
		if test.Severity >= level {
			expected = fmt.Sprintf(noPrefixPattern, test.Severity.Prefix())
		} else {
			expected = "^$"
		}

		s := strings.Title(test.Severity.String())

		buf.Reset()
		test.PrintFunc("hello 23 world")
		validateStd(t, s, expected, buf.String(), test.Severity >= level)

		buf.Reset()
		test.PrintfFunc("hello %d world", 23)
		validateStd(t, fmt.Sprintf("%sf", s), expected, buf.String(), test.Severity >= level)

		buf.Reset()
		test.PrintlnFunc("hello", 23, "world")
		validateStd(t, fmt.Sprintf("%sln", s), expected, buf.String(), test.Severity >= level)
	}

	SetSeverity(severity.Info)
	SetFlags(0)
	SetOutput(os.Stderr)
}

func TestAllStdLogging(t *testing.T) {
	var tests = []struct {
		flag    int
		pattern string // regexp that log output must match; we add `^`, `prefix`, and `expected_text$` always
	}{
		{0, ""},
		{log.Ldate, rStdDate + " "},
		{log.Ltime, rStdTime + " "},
		{log.Ltime | log.Lmicroseconds, rStdTime + rStdMicroseconds + " "},
		{log.Lmicroseconds, rStdTime + rStdMicroseconds + " "}, // microsec implies time
		{log.Llongfile, rStdLongfile + " "},
		{log.Lshortfile, rStdShortfile + " "},
		{log.Llongfile | log.Lshortfile, rStdShortfile + " "}, // shortfile overrides longfile
		{log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile, rStdDate + " " + rStdTime + rStdMicroseconds + " " + rStdLongfile + " "},
		{log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile, rStdDate + " " + rStdTime + rStdMicroseconds + " " + rStdShortfile + " "},
	}

	for _, test := range tests {
		testStdLoggerSeverity(t, test.flag, test.pattern, severity.Debug)
		testStdLoggerSeverity(t, test.flag, test.pattern, severity.Info)
		testStdLoggerSeverity(t, test.flag, test.pattern, severity.Error)
	}
}
