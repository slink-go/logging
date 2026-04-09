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
