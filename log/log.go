/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of maponet/utils.

maponet/utils is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

maponet/utils is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with maponet/utils.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Package log contains a simple multi-level logger.
*/
package log

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

// Available log levels
const (
	ERROR = iota
	INFO
	DEBUG
)

// ErrLevel indicates that the specified log level is not valid
var ErrLevel = errors.New("Unknown log level")

// SetLevel sets the output level for the global logger
func SetLevel(level int) error {
	if level > DEBUG {
		return ErrLevel
	}
	l.level = level
	return nil
}

// SetLevelString sets the output level for the global logger
func SetLevelString(name string) error {
	level, err := StringToLevel(name)
	if err != nil || level > DEBUG {
		return ErrLevel
	}
	l.level = level
	return nil
}

// LevelToString converts a log level to its string representation
func LevelToString(level int) (string, error) {
	var s string
	switch level {
	case ERROR:
		s = "ERROR"
	case INFO:
		s = "INFO"
	case DEBUG:
		s = "DEBUG"
	default:
		return "", ErrLevel
	}
	return s, nil
}

// StringToLevel converts a string to the corresponding log level
func StringToLevel(name string) (int, error) {
	var level int
	switch name {
	case "ERROR":
		level = ERROR
	case "INFO":
		level = INFO
	case "DEBUG":
		level = DEBUG
	default:
		return -1, ErrLevel
	}
	return level, nil
}

func FlagLevel() *string {
	return flag.String("log", "ERROR", "set log level [ERROR|INFO|DEBUG]")
}

type logger struct {
	level int
}

var l logger

func Log(level int, format string, v ...interface{}) {
	if level <= l.level {
		msgType, err := LevelToString(level)
		if err != nil {
			Fatal(err.Error())
		}
		msg := fmt.Sprintf(format, v...)
		t := time.Now().Format(time.RFC1123)
		fmt.Printf("%s [%s]: %s\n", t, msgType, msg)

	}
}

// Fatal logs a message at "ERROR" level and terminates the program
func Fatal(format string, v ...interface{}) {
	Log(ERROR, format, v...)
	os.Exit(1)
}

// Error logs a message at "ERROR" level
func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

// Info logs a message at "INFO" level
func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

// Debug logs a message at "DEBUG" level
func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}
