package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"sync"
	"time"
)

type customLogger struct {
	mu     sync.Mutex
	logFn  func(string)
	logger string
	level  zerolog.Level
}

func newCustomLogger(id string, logFn func(string) /*, opts ...Option*/) Logger {
	return &customLogger{
		level:  getLoggingLevel(id),
		logger: id,
		logFn:  logFn,
	}
}

func (l *customLogger) Clone(newId string) Logger {
	return l
}

func (l *customLogger) Trace(message string, args ...interface{}) {
	if l.IsTraceEnabled() {
		l.log(zerolog.TraceLevel, message, args...)
	}
}
func (l *customLogger) Debug(message string, args ...interface{}) {
	if l.IsDebugEnabled() {
		l.log(zerolog.DebugLevel, message, args...)
	}
}
func (l *customLogger) Info(message string, args ...interface{}) {
	if l.IsInfoEnabled() {
		l.log(zerolog.InfoLevel, message, args...)
	}
}
func (l *customLogger) Warning(message string, args ...interface{}) {
	if l.IsWarningEnabled() {
		l.log(zerolog.WarnLevel, message, args...)
	}
}
func (l *customLogger) Warn(message string, args ...interface{}) {
	l.Warning(message, args...)
}
func (l *customLogger) Error(message string, args ...interface{}) {
	if l.IsErrorEnabled() {
		l.log(zerolog.ErrorLevel, message, args...)
	}
}
func (l *customLogger) Fatal(message string, args ...interface{}) {
	if l.IsFatalEnabled() {
		l.log(zerolog.FatalLevel, message, args...)
	}
}
func (l *customLogger) Panic(message string, args ...interface{}) {
	if l.IsPanicEnabled() {
		l.log(zerolog.PanicLevel, message, args...)
	}
}

func (l *customLogger) IsTraceEnabled() bool {
	return l.level == zerolog.TraceLevel
}
func (l *customLogger) IsDebugEnabled() bool {
	return l.level <= zerolog.DebugLevel
}
func (l *customLogger) IsInfoEnabled() bool {
	return l.level <= zerolog.InfoLevel
}
func (l *customLogger) IsWarningEnabled() bool {
	return l.level <= zerolog.WarnLevel
}
func (l *customLogger) IsErrorEnabled() bool {
	return l.level <= zerolog.ErrorLevel
}
func (l *customLogger) IsFatalEnabled() bool {
	return l.level <= zerolog.FatalLevel
}
func (l *customLogger) IsPanicEnabled() bool {
	return l.level <= zerolog.PanicLevel
}

func (l *customLogger) SetLevel(level string) Logger {
	var err error
	l.level, err = stringToLevel(level)
	if err != nil {
		l.level = zerolog.InfoLevel
	}
	return l
}
func (l *customLogger) GetLevel() string {
	return l.level.String()
}

func (l *customLogger) log(level zerolog.Level, message string, args ...interface{}) {
	if l.logFn == nil {
		return
	}
	msg := l.format(level, message, args...)
	l.mu.Lock()
	l.logFn(msg)
	l.mu.Unlock()
}
func (l *customLogger) format(level zerolog.Level, message string, args ...interface{}) string {
	return fmt.Sprintf(
		fileLogFormat,
		time.Now().Format("2006-01-02 15:04:05.000"),
		logLevelAbbr(level),
		l.logger,
		message,
		l.args(args...),
	)
}
func (l *customLogger) args(args ...interface{}) string {
	var sb strings.Builder
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			sb.WriteString(fmt.Sprintf(" %v=%v", args[i], args[i+1]))
		}
	}
	return sb.String()
}
