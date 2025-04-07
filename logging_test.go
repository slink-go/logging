package logging

import "testing"

func TestLoggerWithCaller(t *testing.T) {

	GetLogger("test1").Info("test message 1.1")
	GetLogger("test1").Info("test message 1.2")
	GetLogger("test2", WithCaller()).Info("test message 2.1")
	GetLogger("test2").Info("test message 2.2")

}
