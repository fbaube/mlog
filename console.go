// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	S "strings"
)

type consoleBrush func(string) string

func newConsoleBrush(format string) consoleBrush {
	return func(text string) string {
		return "\033[" + format + "m" + text + "\033[0m"
	}
}

/*
       FG  BG
Black: 30  40
Red:   31  41
Green: 32  42
Yello: 33  43
Blue:  34  44
Mgnta: 35  45
Cyan:  36  46
White: 37  47
Reset: 0 (all)

Bold   ;1
Dim    ;2
Italic ;3
Undrln ;4

Rvrsd  ;7
Strkthru ;9
*/

var Brushes = map[Level]consoleBrush{
	LevelDbg:      newConsoleBrush("1;36"),   // bold cyan
	LevelProgress: newConsoleBrush("36"),     // cyan
	LevelInfo:     newConsoleBrush("36"),     // cyan
	LevelOkay:     newConsoleBrush("32"),     // green
	LevelWarning:  newConsoleBrush("33"),     // yellow
	LevelError:    newConsoleBrush("31"),     // red
	LevelPanic:    newConsoleBrush("1;95"),   // bold light magenta
	GreenBG:       newConsoleBrush("42;2;4"), // green background
}

// ConsoleTarget writes filtered log messages to console window.
type ConsoleTarget struct {
	*Filter
	ColorMode   bool      // whether to use colors to differentiate log levels
	Writer      io.Writer // the writer to write log messages
	close       chan bool
	DetailsInfo        // NEW
	Category    string // NEW
	Subcategory string // NEW
}

func (t *ConsoleTarget) SetCategory(s string) {
	t.Category = s
}

func (t *ConsoleTarget) SetSubcategory(s string) {
	t.Subcategory = s
}

// NewConsoleTarget creates a ConsoleTarget.
// The new ConsoleTarget takes these default options:
// MaxLevel: LevelDebug, ColorMode: true, Writer: os.Stdout
func NewConsoleTarget() *ConsoleTarget {
	return &ConsoleTarget{
		Filter:    &Filter{MaxLevel: LevelDbg},
		ColorMode: true,
		Writer:    os.Stdout,
		close:     make(chan bool, 0),
	}
}

// Open prepares ConsoleTarget for processing log messages.
func (t *ConsoleTarget) Open(io.Writer) error {
	t.Filter.Init()
	if t.Writer == nil {
		return errors.New("ConsoleTarget.Writer cannot be nil")
	}
	if runtime.GOOS == "windows" {
		t.ColorMode = false
	}
	return nil
}

// Process writes a log message using Writer.
func (t *ConsoleTarget) Process(e *Entry) {
	if e == nil {
		t.close <- true
		return
	}
	if !t.Allow(e) {
		return
	}
	msg := e.String()
	if t.ColorMode {
		if !S.Contains(msg, "\033[") {
			brush, ok := Brushes[e.Level]
			if ok {
				msg = brush(msg)
			}
		}
	}
	fmt.Fprintln(t.Writer, msg)
}

// Close closes the console target.
func (t *ConsoleTarget) Close() {
	<-t.close
}

// Flush is a no-op.
func (t *ConsoleTarget) Flush() {
}

func (t *ConsoleTarget) DoesDetails() bool {
	return true
}

func (t *ConsoleTarget) StartDetailsBlock(*Entry) {
	fmt.Fprintln(t.Writer, "NOT IMPLEMENTED YET: ConsoleTarget.StartDetailsBlock")
}

func (t *ConsoleTarget) CloseDetailsBlock(string) {
	fmt.Fprintln(t.Writer, "NOT IMPLEMENTED YET: ConsoleTarget.CloseDetailsBlock")
}
