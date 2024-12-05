package logging

import (
	"github.com/rs/zerolog"
)

// region - file logger

type noOpLogger struct {
}

func (l *noOpLogger) Clone(newId string) Logger {
	return l
}

func (l *noOpLogger) Trace(message string, args ...interface{}) {
}
func (l *noOpLogger) Debug(message string, args ...interface{}) {
}
func (l *noOpLogger) Info(message string, args ...interface{}) {
}
func (l *noOpLogger) Warning(message string, args ...interface{}) {
}
func (l *noOpLogger) Error(message string, args ...interface{}) {
}
func (l *noOpLogger) Fatal(message string, args ...interface{}) {
}
func (l *noOpLogger) Panic(message string, args ...interface{}) {
}

func (l *noOpLogger) IsTraceEnabled() bool {
	return false
}
func (l *noOpLogger) IsDebugEnabled() bool {
	return false
}
func (l *noOpLogger) IsInfoEnabled() bool {
	return false
}
func (l *noOpLogger) IsWarningEnabled() bool {
	return false
}
func (l *noOpLogger) IsErrorEnabled() bool {
	return false
}
func (l *noOpLogger) IsFatalEnabled() bool {
	return false
}
func (l *noOpLogger) IsPanicEnabled() bool {
	return false
}

func (l *noOpLogger) SetLevel(level string) {
}
func (l *noOpLogger) GetLevel() string {
	return zerolog.Disabled.String()
}

// endregion
