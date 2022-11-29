package log_test

import (
	L "github.com/fbaube/mlog"
)

func Example() {
	// MAYBE not necessary! Does running an Example func invoke L.init() ?
	LL := L.NewLogger()
	// Add two loggers that write to Stdout
	LL.Targets = append(L.L.Targets, L.NewConsoleTarget())
	LL.Targets = append(L.L.Targets, L.NewConsoleTarget())
	LL.Open()
	LL.Info("%d", len(L.L.Targets))
	LL.Dbg("Debug message")
	LL.Error("Error message")
	// Output:
	// 2
	// Debug message
	// Error message
}
