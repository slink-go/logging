package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

func TestLoggerBasic(t *testing.T) {
	os.Setenv("LOGGING_FORMAT", "pretty")
	GetLogger("test").SetLevel("TRACE").Info("message with args", "k1", "v", "k2", 123, "k3", 0.1, "k4", false)
	GetLogger("test").SetLevel("TRACE").Info("message basic")
}

func TestLoggerWith(t *testing.T) {
	os.Setenv("LOGGING_FORMAT", "pretty")
	GetLogger("test", With("key", "value", "num", 10)).SetLevel("TRACE").Info("message with args", "k1", "v", "k2", 123, "k3", 0.1, "k4", false)
	GetLogger("test", With("key", "value", "num", 10)).Info("message basic")
}

func TestLoggerWithCaller(t *testing.T) {

	os.Setenv("LOGGING_FORMAT", "pretty")

	GetLogger("testY").Info("test message")
	GetLogger("testX", WithCaller()).Info("test message; skip = default (3)")
	GetLogger("test0", WithCaller(0)).Info("test message; skip = 0")
	GetLogger("test1", WithCaller(1)).Info("test message; skip = 1")
	GetLogger("test2", WithCaller(2)).Info("test message; skip = 2")
	GetLogger("test3", WithCaller(3)).Info("test message; skip = 3")
	GetLogger("test4", WithCaller(4)).Info("test message; skip = 4")
	GetLogger("test5", WithCaller(5)).Info("test message; skip = 5")

	testLogLevel1(GetLogger("test3"))
}

type testHook struct {
}

func (t testHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	fmt.Println("hook message:", message)
	fmt.Println("hook level:", level)
	fmt.Printf("hook event: %#v\n", *e)
}
func TestLoggerWithHook(t *testing.T) {
	GetLogger("testY", WithHook(testHook{})).Info("test message")
}
func TestLoggerWithLevel(t *testing.T) {
	os.Setenv("LOGGING_FORMAT", "pretty")
	os.Setenv("LOGGING_LEVEL_ROOT", "info")
	GetLogger("testLvl1", WithLevelStr("TRACE")).Debug("test message debug")
}

func testLogLevel1(logger Logger) {
	logger.Info("test message level 1")
	testLogLevel2(logger)
}
func testLogLevel2(logger Logger) {
	logger.Info("test message level 2")
	testLogLevel3(logger)
}
func testLogLevel3(logger Logger) {
	logger.Info("test message level 3")
}
