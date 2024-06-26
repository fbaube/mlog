// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"fmt"
	LU "github.com/fbaube/logutils"
	"io"
	"os"
)

// FileTarget writes filtered log messages to a file.
// FileTarget supports file rotation by keeping certain number of backup log files.
type FileTarget struct {
	*Filter
	// the log file name. When Rotate is true, log file name will be suffixed
	// to differentiate different backup copies (e.g. app.log.1)
	FileName string
	// whether to enable file rotating at specific time interval or when maximum file size is reached.
	Rotate bool
	// how many log files should be kept when Rotate is true (the current log file is not included).
	// This field is ignored when Rotate is false.
	BackupCount int
	// maximum number of bytes allowed for a log file. Zero means no limit.
	// This field is ignored when Rotate is false.
	MaxBytes int64

	fd           *os.File
	currentBytes int64
	errWriter    io.Writer
	close        chan bool

	Category    string
	Subcategory string
}

func (t *FileTarget) SetCategory(s string) {
	t.Category = s
}

func (t *FileTarget) SetSubcategory(s string) {
	t.Subcategory = s
}

// NewFileTarget creates a FileTarget.
// The new FileTarget takes these default options:
// MaxLevel: LevelNotice, Rotate: true, BackupCount: 10, MaxBytes: 1 << 20
// After calling this, you must fill in the FileName field.
func NewFileTarget() *FileTarget {
	return &FileTarget{
		Filter:      &Filter{MaxLevel: LU.LevelInfo},
		Rotate:      true,
		BackupCount: 10,
		MaxBytes:    1 << 20, // 1MB
		close:       make(chan bool, 0),
	}
}

// Open prepares FileTarget for processing log messages.
func (t *FileTarget) Open(errWriter io.Writer) error {
	t.Filter.Init()
	if t.FileName == "" {
		return errors.New("FileTarget.FileName must be set")
	}
	if t.Rotate {
		if t.BackupCount < 0 {
			return errors.New("FileTarget.BackupCount must be no less than 0")
		}
		if t.MaxBytes <= 0 {
			return errors.New("FileTarget.MaxBytes must be no less than 0")
		}
	}

	fd, err := os.OpenFile(t.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("FileTarget was unable to create a log file: %v", err)
	}
	t.fd = fd
	t.errWriter = errWriter

	return nil
}

// Process saves an allowed log message into the log file.
func (t *FileTarget) Process(e *Entry) {
	if e == nil {
		t.fd.Close()
		t.close <- true
		return
	}
	if t.fd != nil && t.Allow(e) {
		if t.Rotate {
			t.rotate(int64(len(e.String()) + 1))
		}
		n, err := t.fd.Write([]byte(e.String() + "\n"))
		t.currentBytes += int64(n)
		if err != nil {
			fmt.Fprintf(t.errWriter, "FileTarge write error: %v\n", err)
		}
	}
}

// Close closes the file target.
func (t *FileTarget) Close() {
	<-t.close
}

// Flush is a no-op but SHOULD have a call to flush the file.
func (t *FileTarget) Flush() {
	// t.fd.Flush()
}

func (t *FileTarget) DoesDetails() bool {
	return true
}

func (t *FileTarget) StartDetailsBlock(*Entry) {
	fmt.Fprintln(t.fd, "NOT IMPLEMENTED YET: FileTarget.StartDetailsBlock")
}

func (t *FileTarget) CloseDetailsBlock(string) {
	fmt.Fprintln(t.fd, "NOT IMPLEMENTED YET: FileTarget.CloseDetailsBlock")
}

func (t *FileTarget) rotate(bytes int64) {
	if t.currentBytes+bytes <= t.MaxBytes || bytes > t.MaxBytes {
		return
	}
	t.fd.Close()
	t.currentBytes = 0

	var err error
	for i := t.BackupCount; i >= 0; i-- {
		path := t.FileName
		if i > 0 {
			path = fmt.Sprintf("%v.%v", t.FileName, i)
		}
		if _, err = os.Lstat(path); err != nil {
			// file not exists
			continue
		}
		if i == t.BackupCount {
			os.Remove(path)
		} else {
			os.Rename(path, fmt.Sprintf("%v.%v", t.FileName, i+1))
		}
	}
	t.fd, err = os.OpenFile(t.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		t.fd = nil
		fmt.Fprintf(t.errWriter, "FileTarget was unable to create a log file: %v", err)
	}
}
