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
	"fmt"
	"io"
	"os"
	"time"
)

// Available log levels
const (
	ERROR = iota
	INFO
	DEBUG
)

// BadLevelError indicates that the specified log level is not valid.
type BadLevelError struct {
	level interface{}
}

func (e *BadLevelError) Error() string { return fmt.Sprintf("Unknown log level: %v", e.level) }

// Logger is a simple multi-level logger.
type Logger struct {
	Level  int
	Writer io.Writer
}

func NewLogger(level interface{}, writer io.Writer) *Logger {
	var l Logger
	l.SetLevel(level)
	l.Writer = writer
	return &l
}

// SetLevel sets the output level for the logger.
// 'level' can be either a string or an int.
func (l *Logger) SetLevel(level interface{}) error {
	switch level {
	case ERROR, "ERROR":
		l.Level = ERROR
	case INFO, "INFO":
		l.Level = INFO
	case DEBUG, "DEBUG":
		l.Level = DEBUG
	default:
		return &BadLevelError{level}
	}
	return nil
}

// Log logs a message with custom level and type.
// To log messages at predefined levels you can use the convenience
// functions Fatal(), Error(), Info(), Debug().
func (l *Logger) Log(level int, msgType string, format string, v ...interface{}) {
	if level <= l.Level {
		msg := fmt.Sprintf(format, v...)
		t := time.Now().Format(time.RFC1123)
		fmt.Fprintf(l.Writer, "%s [%s]: %s\n", t, msgType, msg)

	}
}

// Fatal logs a message at "ERROR" level and terminates the program.
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Log(ERROR, "ERROR", format, v...)
	os.Exit(1)
}

// Error logs a message at "ERROR" level.
func (l *Logger) Error(format string, v ...interface{}) {
	l.Log(ERROR, "ERROR", format, v...)
}

// Info logs a message at "INFO" level.
func (l *Logger) Info(format string, v ...interface{}) {
	l.Log(INFO, "INFO", format, v...)
}

// Debug logs a message at "DEBUG" level.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.Log(DEBUG, "DEBUG", format, v...)
}

var DefaultLogger = NewLogger(ERROR, os.Stdout)

// SetLevel sets the output level for the default logger.
func SetLevel(level interface{}) error {
	return DefaultLogger.SetLevel(level)
}

// Fatal logs a message at "ERROR" level with the default logger and
// terminates the program.
func Fatal(msg string, v ...interface{}) {
	DefaultLogger.Fatal(msg, v...)
}

// Error logs a message at "ERROR" level with the default logger.
func Error(msg string, v ...interface{}) {
	DefaultLogger.Error(msg, v...)
}

// Info logs a message at "INFO" level with the default logger.
func Info(msg string, v ...interface{}) {
	DefaultLogger.Info(msg, v...)
}

// Debug logs a message at "DEBUG" level with the default logger.
func Debug(msg string, v ...interface{}) {
	DefaultLogger.Debug(msg, v...)
}
