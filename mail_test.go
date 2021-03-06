// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log_test

import (
	"testing"

	log "github.com/fbaube/mlog"
)

func TestNewMailTarget(t *testing.T) {
	target := log.NewMailTarget()
	if target.MaxLevel != log.LevelDbg {
		t.Errorf("NewMailTarget.MaxLevel = %v, expected %v", target.MaxLevel, log.LevelDbg)
	}
}
