package logging

import (
	"fmt"
	"testing"
)

func logFn(message string) {
	fmt.Println(message)
}

func TestCustomLogger(t *testing.T) {
	//os.Setenv("LOGGING_LEVEL_ROOT", "trace")
	l := GetCustomLogger("test", logFn)
	l.SetLevel("trace")
	l.Trace("hello trace")
	l.Debug("hello debug", "a", 1, "b", "2")
	l.Info("hello info")
	l.Warning("hello warn", "a", 1, "b", "2")
	l.Error("hello error")
	l.Clone("another").Panic("hello panic", "a", 1, "b", "2")
}

func TestCustomLoggerTs(t *testing.T) {
	l1 := GetCustomLogger("l1", logFn)
	l2 := GetCustomLoggerWithTimestamp("l2", logFn)
	l1.Info("hello l1")
	l2.Info("hello l2")
}
