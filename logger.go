// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package log implements logging with severity levels and message categories.
package log

// TODO: If "app" is default, can be omitted from log msgs.
// if is category, and filterable, then use for go pkg name, or for input file name.

import (
	"bytes"
	"errors"
	"fmt"
	LU "github.com/fbaube/logutils"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level describes the level of a log message.
// type Level int

// L is the predefined default global logger.
var L *Logger

// pCoreLogger is a cheat so that we can set max logging level globally.
var pCoreLogger *Logger

// init calls NewLogger() to create a root logger.
// The root logger is os.Stderr, LevelDebug.
func init() {
	L = NewLogger()
}

// Entry represents a log entry.
type Entry struct {
	Level            LU.Level
	Category         string
	Message          string
	Time             time.Time
	CallStack        string
	FormattedMessage string
}

// String returns the string representation of the log entry
func (e *Entry) String() string {
	return e.FormattedMessage
}

// Target represents a target where the logger can
// send log messages to for further processing.
type Target interface {
	// Open prepares the target for processing log messages.
	// Called when Logger.Open() is called.
	// If an error is returned, the target will be removed from the
	// logger. errWriter should be used to write errors found while
	// processing log messages, and should probably default to Stderr.
	Open(errWriter io.Writer) error
	// Process processes an incoming log message.
	Process(*Entry)
	// Close closes a target.
	// Called when Logger.Close() is called. Each target gets
	// a chance to flush its log messages to its destination.
	Close()
	// Flush is NEW and added so that logging plays nicely with
	// other sources of text.
	Flush()
	// DoesDetails is NEW and has a value per-struct, not per-instance.
	DoesDetails() bool
}

// coreLogger maintains the log messages in a channel and sends them to various targets.
type coreLogger struct {
	lock        sync.Mutex
	open        bool        // whether the logger is open
	entries     chan *Entry // log entries
	ErrorWriter io.Writer   // the writer to record errors caused by log targets

	BufferSize     int // the size of the channel storing log entries
	CallStackDepth int // the number of call stack frames to log for each
	//                 // message. 0 means do not log any call stack frame.
	CallStackFilter string // a substring that a call stack frame filepath
	//                     // should contain in order for the frame to be counted
	MaxLevel LU.Level // the maximum level of messages to be logged
	Targets  []Target // targets for sending log messages to
}

// Formatter formats a log message into an appropriate string.
type Formatter func(*Logger, *Entry) string

// Logger records log messages and dispatches them to various targets for further processing.
type Logger struct {
	*coreLogger
	Category  string    // the category associated with this logger
	Formatter Formatter // message formatter
}

// NewLogger creates a root logger.
// The new logger takes these default options:
// ErrorWriter: os.Stderr, BufferSize: 1024, MaxLevel: LU.LevelDebug,
// Category: app, Formatter: DefaultFormatter
func NewLogger() *Logger {
	logger := &coreLogger{
		ErrorWriter: os.Stderr,
		BufferSize:  1024,
		MaxLevel:    LU.LevelDebug,
		Targets:     make([]Target, 0),
	}
	pCoreLogger = &Logger{logger, "", DefaultFormatter}
	return pCoreLogger // &Logger{logger, "", DefaultFormatter}
}

// NewNullLogger creates a no-op logger.
// .
func NewNullLogger() *Logger {
	logger := &coreLogger{
		ErrorWriter: io.Discard,
		BufferSize:  1024,
		MaxLevel:    LU.LevelError,
		Targets:     make([]Target, 0),
	}
	pCoreLogger = &Logger{logger, "", DefaultFormatter}
	return pCoreLogger // &Logger{logger, "", DefaultFormatter}
}

// GetLogger creates a logger with the specified category and log formatter.
// Messages logged thru this logger will carry the same category name.
// The formatter, if not specified, will inherit from the calling logger.
// It will be used to format all messages logged thru this logger.
func (l *Logger) GetLogger(category string, formatter ...Formatter) *Logger {
	if len(formatter) > 0 {
		return &Logger{l.coreLogger, category, formatter[0]}
	}
	return &Logger{l.coreLogger, category, l.Formatter}
}

// Panic logs a message indicating the system is dying,
// but does NOT actually execute a call to panic(..)
func (l *Logger) Panic(format string, a ...interface{}) {
	l.Log(LU.LevelPanic, format, a...)
}

// Error logs a message indicating an error condition.
// This method takes one or multiple parameters. If a
// single parameter is provided, it IS the log message.
// If multiple parameters are provided, they are passed
// to fmt.Sprintf() to generate the log message.
func (l *Logger) Error(format string, a ...interface{}) {
	l.Log(LU.LevelError, format, a...)
}

// Warning logs a message indicating a warning condition.
func (l *Logger) Warning(format string, a ...interface{}) {
	l.Log(LU.LevelWarning, format, a...)
}

// Okay logs a message indicating an okay condition.
func (l *Logger) Okay(format string, a ...interface{}) {
	l.Log(LU.LevelOkay, format, a...)
}

// Info logs a message for a normal but meaningful condition.
func (l *Logger) Info(format string, a ...interface{}) {
	l.Log(LU.LevelInfo, format, a...)
}

/*
// Progress logs a message for how things are progressing.
func (l *Logger) Progress(format string, a ...interface{}) {
	l.Log(LU.LevelProgress, format, a...)
}
*/

// Debug logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debug(format string, a ...interface{}) {
	l.Log(LU.LevelDebug, format, a...)
}

// Log logs a message of a specified severity level.
func (l *Logger) Log(level LU.Level, format string, a ...interface{}) {
	if level > l.MaxLevel || !l.open {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	entry := &Entry{
		Category: l.Category,
		Level:    level,
		Message:  message,
		Time:     time.Now(),
	}
	if l.CallStackDepth > 0 {
		entry.CallStack = GetCallStack(3, l.CallStackDepth, l.CallStackFilter)
	}
	entry.FormattedMessage = l.Formatter(l, entry)
	l.entries <- entry
}

func (l *Logger) LogWithString(level LU.Level, format string, special string, a ...interface{}) {
	// func (l *Logger) Log(level LU.Level, format string, a ...interface{}) {
	if level > l.MaxLevel || !l.open {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	entry := &Entry{
		Category: l.Category,
		Level:    level,
		Message:  "(" + special + ") " + message,
		Time:     time.Now(),
	}
	if l.CallStackDepth > 0 {
		entry.CallStack = GetCallStack(3, l.CallStackDepth, l.CallStackFilter)
	}
	entry.FormattedMessage = l.Formatter(l, entry)
	l.entries <- entry
}

func SetMaxLevel(lvl LU.Level) {
	pCoreLogger.MaxLevel = lvl
}

// Open prepares the logger and the targets for logging purpose.
// Open must be called before any message can be logged.
func (l *coreLogger) Open() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.open {
		return nil
	}
	if l.ErrorWriter == nil {
		return errors.New("Logger.ErrorWriter must be set.")
	}
	if l.BufferSize < 0 {
		return errors.New("Logger.BufferSize must be no less than 0.")
	}
	if l.CallStackDepth < 0 {
		return errors.New("Logger.CallStackDepth must be no less than 0.")
	}
	l.entries = make(chan *Entry, l.BufferSize)
	var targets []Target
	for _, target := range l.Targets {
		if err := target.Open(l.ErrorWriter); err != nil {
			fmt.Fprintf(l.ErrorWriter, "Failed to open target: %v", err)
		} else {
			targets = append(targets, target)
		}
	}
	l.Targets = targets
	go l.process()
	l.open = true
	return nil
}

// process sends the messages to targets for processing.
func (l *coreLogger) process() {
	for {
		entry := <-l.entries
		for _, target := range l.Targets {
			target.Process(entry)
		}
		if entry == nil {
			break
		}
	}
}

// Close closes the logger and the targets.
// Existing messages will be processed before the targets are closed.
// New incoming messages will be discarded after calling this method.
func (l *coreLogger) Close() {
	if !l.open {
		return
	}
	l.open = false
	// use a nil entry to signal the close of logger
	l.entries <- nil
	for _, target := range l.Targets {
		target.Close()
	}
}

// Flush flushes the logger and the targets.
func (l *coreLogger) Flush() {
	if !l.open {
		return
	}
	for _, target := range l.Targets {
		target.Flush()
	}
}

// DefaultFormatter is the default formatter used to format every log message.
// This formatter assumes no Target is a DetailsTarget.
func DefaultFormatter(l *Logger, e *Entry) string {
	var sTime, sLvl, sCtg string
	sLvl = e.Level.String()
	if len(sLvl) != 5 {
		sLvl = sLvl[0:4]
	}
	sTime = e.Time.Format("15.04.05") // e.Time.Format("01-02-15.04.05")
	if e.Category != "" {
		sCtg = fmt.Sprintf("[%s]", e.Category)
	}
	return fmt.Sprintf("%s %s"+ /*[%s]*/ "%s %v %v",
		sTime, LU.EmojiOfLevel(e.Level), // sLvl,
		sCtg, e.Message, e.CallStack)
}

// GetCallStack returns the current call stack information as a string.
// The skip parameter specifies how many top frames should be skipped, while
// the frames parameter specifies at most how many frames should be returned.
func GetCallStack(skip int, frames int, filter string) string {
	buf := new(bytes.Buffer)
	for i, count := skip, 0; count < frames; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if filter == "" || strings.Contains(file, filter) {
			fmt.Fprintf(buf, "\n%s:%d", file, line)
			count++
		}
	}
	return buf.String()
}
