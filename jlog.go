/*
Copyright (c) 2018-2019 Drew DeVault
Copyright (c) 2021-2022 Robin Jarry

The MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package jlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	dbg    *log.Logger
	warn   *log.Logger
	info   *log.Logger
	errorl *log.Logger
)

const (
	fileLog         = "jasmd.log"
	clientsLogFile  = "jasmd_clients.log"
	jasmDir         = "jasm"
	xdgStateDefault = ".local/state"
	xdgStateEnv     = "XDG_STATE_HOME"
	//Error bool
)

// edit
var (
	Debug bool
)

func CloseLog(f *os.File) {
	f.Close()
}

func setStateDir() (string, error) {
	var (
		stateDir string
		ok       bool
	)
	stateDir, ok = os.LookupEnv(xdgStateEnv)
	if !ok {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		stateDir = filepath.Join(homeDir, xdgStateDefault)
	}
	jasmStateDir := filepath.Join(stateDir, jasmDir)
	err := os.MkdirAll(jasmStateDir, 0700)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("%v", err)
	}
	return jasmStateDir, nil
}

//https://git.sr.ht/~rjarry/aerc/tree/master/item/log/logger.go
// https://www.honeybadger.io/blog/golang-logging/
// https://stackoverflow.com/questions/36719525/how-to-log-messages-to-the-console-and-a-file-both-in-golang
func InitLog() (*os.File, error) {
	stateDir, err := setStateDir()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	f := filepath.Join(stateDir, fileLog)

	logFile, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("open logFile: %v", err)
	}
	mw := io.MultiWriter(os.Stderr, logFile)

	info = log.New(mw, "[jasmd] ", 0)
	warn = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile) // | bitwise OR
	dbg = log.New(mw, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)    // | bitwise OR
	//errorl = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorl = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return logFile, nil
}

func InitClientLog() (*os.File, error) {
	stateDir, err := setStateDir()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	f := filepath.Join(stateDir, clientsLogFile)

	logFile, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("open nsmclients logFile: %v", err)
	}
	return logFile, nil
}

func ErrorLogger() *log.Logger {
	if errorl == nil {
		return log.New(io.Discard, "", log.LstdFlags)
	}
	return errorl
}

func Debugf(message string, args ...interface{}) {
	if dbg == nil || !Debug { // NOTE edit
		return
	}
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	dbg.Output(2, message) //nolint:errcheck // we can't do anything with what we log
}

func Infof(msg string, args ...interface{}) {
	if info == nil { // NOTE always info
		return
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	info.Output(1, msg) //nolint:errcheck // we can't do anything with what we log
}

/*

func Output(calldepth int, s string) error
Calldepth is the count of the number of frames to skip when computing the file name and line number if Llongfile or Lshortfile is set; a value of 1 will print the details for the caller of Output.

*/

func Warnf(message string, args ...interface{}) {
	if warn == nil { // NOTE always warn.
		return
	}
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	warn.Output(2, message) //nolint:errcheck // we can't do anything with what we log
}

func Errorf(message string, args ...interface{}) {
	if errorl == nil { // || !Error {
		return
	}
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	errorl.Output(2, message) //nolint:errcheck // we can't do anything with what we log
}
