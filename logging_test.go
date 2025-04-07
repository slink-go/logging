package logging

import (
	"testing"
)

func TestLoggerWithCaller(t *testing.T) {

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
