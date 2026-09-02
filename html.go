package log

import "io"

type HtmlTarget struct {
	*Filter
	// the target HTML element's ID attribute.
	FieldID   string
	Writer    io.Writer // the writer to write log messages
	errWriter io.Writer
	close     chan bool
	DetailsInfo
}
