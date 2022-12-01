package log

import (
	"fmt"
	LU "github.com/fbaube/logutils"
	S "strings"
)

// DetailsFormatter formats a log message into an appropriate string,
// but also handles Details nicely (in a manner TBD), and accepts a
// third argument of a string that may be specified per-message.
type DetailsFormatter func(*Logger, *Entry, []string) string

// DetailsInfo is embedded in details-capable Target's. It
// applies only to "log details", which logging is stateful,
// not to "text quotes", which logging is an atomic operation.
type DetailsInfo struct {
	DoingDetails     bool
	MinLogLevel      LU.Level
	Category         string
	Subcategory      string
	DetailsFormatter // message formatter
}

// DetailsTarget is a target where the logger can both
// (1) Open a collapsible, ignorable set of log detail messages, and
// (2) Quote a collapsible, ignorable block of text (atomic operation!).
//
// In a Console target, do this by omitting the first three (or
// six) characters of the timestamp, providing visual indenting.
// For (1) use " - " or " * ", so that it resembles a list.
// For (2) use " " " or " ' ", so that it is obv a quote.
//
// In an HTML target, do this by opening a "<details> block" and
// in the very same log message, opening the <summary>  element.
// Then subsequent log messages (or the accompanying text block)
// can be written to the body of the <details> element (separated
// by <br/> tags, rather than by newlines as in most log targets)
// until the <details> element is closed.
//
// As an enhancement, a set of log detail messages tracks its minimum
// (i.e. most severe) logging level, with a summary line at the end.
//
// The five function calls could be ignored as no-ops by targets
// that do not implemement the interface. However it is simpler
// and clearer, and follows the existing processing architecture,
// to make the call for every Target in a Logger, but then in the
// Logger method call, check each Target for being a DetailsTarget.
type DetailsTarget interface {
	Target
	StartLogDetailsBlock(string, *Entry) // s = Category e.g. "[01]" and clear Subcat
	CloseLogDetailsBlock(string)
	LogTextQuote(*Entry, string)
	// SetCategory is used per-Contentity, e.g. "[00]", "[01]", ...
	SetCategory(string)
	// SetSubcategory is used per- Contentity processing stage, e.g. "[st1b]"
	SetSubcategory(string)
}

/*
StartLogDetailsBlock(string, *Entry) // s = Category e.g. "[01]" and clear Subcat
CloseLogDetailsBlock(string)
*/

// SetCategory is for DetailsTarget's.
func (l *coreLogger) SetCategory(s string) {
	if !l.open {
		return
	}
	for _, target := range l.Targets {
		dt, OK := target.(DetailsTarget)
		if OK {
			dt.SetCategory(s)
		}
	}
}

// SetSubcategory is for DetailsTarget's.
func (l *coreLogger) SetSubcategory(s string) {
	if !l.open {
		return
	}
	for _, target := range l.Targets {
		dt, OK := target.(DetailsTarget)
		if OK {
			dt.SetSubcategory(s)
		}
	}
}

// DefaultDetailsFormatter is the default formatter used to format every
// log message when the Target is details-capable. In this formatter, we
// assume that the Logger IS a Details Logger.
//
// Note that this only really works with single threading, or else the
// log messages of different Details sets get all mised up.
func DefaultDetailsFormatter(l *Logger, e *Entry, spcl []string) string {
	var sTime, sLvl, sCtg, sSpcl string
	sLvl = e.Level.String()
	if len(sLvl) != 5 {
		sLvl = sLvl[0:4]
	}
	sTime = e.Time.Format("15.04.05") // e.Time.Format("01-02-15.04.05")
	if e.Category != "" {
		sCtg = "[" + e.Category + "]"
	}
	if spcl != nil {
		var sb S.Builder
		sb.Reset()
		sb.Grow(20)
		for _, ss := range spcl {
			sb.WriteString("," + ss)
		}
		sSpcl = " (" + sb.String()[1:] + ") "
	}
	return fmt.Sprintf("%s %s%s[%s]%s %v %v",
		sTime, sSpcl, LU.EmojiOfLevel(e.Level), sLvl, sCtg,
		e.Message, e.CallStack)
}

func LogTextQuote(*Entry, string) {
	panic("LogTextQuote")
}
