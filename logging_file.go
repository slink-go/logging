package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// region - file logger

var fileMutex map[string]*sync.Mutex

func getMutex(file *os.File) *sync.Mutex {
	if file == nil {
		return nil
	}
	if fileMutex == nil {
		fileMutex = make(map[string]*sync.Mutex)
	}
	ap, err := filepath.Abs(file.Name())
	if err != nil {
		return nil
	}
	v, ok := fileMutex[ap]
	if !ok {
		v = &sync.Mutex{}
		fileMutex[ap] = v
		return v
	}
	return v
}

func newCommonFileLogger(f *os.File, id string) Logger {
	return &fileLogger{
		level:  getLoggingLevel(id),
		logger: id,
		file:   f,
	}
}

const fileLogFormat = "%s %3s [%s] %s"

type fileLogger struct {
	//mutex  sync.Mutex
	level  zerolog.Level
	file   *os.File
	logger string
}

func (l *fileLogger) log(level zerolog.Level, message string, args ...interface{}) {
	if l.file == nil {
		return
	}
	msg := l.format(level, message, args...)
	mtx := getMutex(l.file)
	if mtx == nil {
		return
	}
	mtx.Lock()
	_, err := fmt.Fprintln(l.file, msg)
	if err != nil {
		panic(err)
	}
	mtx.Unlock()
}
func (l *fileLogger) format(level zerolog.Level, message string, args ...interface{}) string {
	return fmt.Sprintf(
		fileLogFormat,
		time.Now().Format("2006-01-02 15:04:05.000"),
		logLevelAbbr(level),
		l.logger,
		fmt.Sprintf(message, args...),
	)
}

func (l *fileLogger) Clone(newId string) Logger {
	return GetFileLogger(l.file, newId)
}

func (l *fileLogger) Trace(message string, args ...interface{}) {
	if l.level == zerolog.TraceLevel {
		l.log(zerolog.TraceLevel, message, args...)
	}
}
func (l *fileLogger) Debug(message string, args ...interface{}) {
	if l.level <= zerolog.DebugLevel {
		l.log(zerolog.DebugLevel, message, args...)
	}
}
func (l *fileLogger) Info(message string, args ...interface{}) {
	if l.level <= zerolog.InfoLevel {
		l.log(zerolog.InfoLevel, message, args...)
	}
}
func (l *fileLogger) Warning(message string, args ...interface{}) {
	if l.level <= zerolog.WarnLevel {
		l.log(zerolog.WarnLevel, message, args...)
	}
}
func (l *fileLogger) Error(message string, args ...interface{}) {
	if l.level <= zerolog.ErrorLevel {
		l.log(zerolog.ErrorLevel, message, args...)
	}
}
func (l *fileLogger) Fatal(message string, args ...interface{}) {
	if l.level <= zerolog.FatalLevel {
		l.log(zerolog.FatalLevel, message, args...)
	}
}
func (l *fileLogger) Panic(message string, args ...interface{}) {
	if l.level <= zerolog.PanicLevel {
		l.log(zerolog.PanicLevel, message, args...)
		panic(fmt.Sprintf(message, args...))
	}
}

func (l *fileLogger) IsTraceEnabled() bool {
	return l.level == zerolog.TraceLevel
}
func (l *fileLogger) IsDebugEnabled() bool {
	return l.level <= zerolog.DebugLevel
}
func (l *fileLogger) IsInfoEnabled() bool {
	return l.level <= zerolog.InfoLevel
}
func (l *fileLogger) IsWarningEnabled() bool {
	return l.level <= zerolog.WarnLevel
}
func (l *fileLogger) IsErrorEnabled() bool {
	return l.level <= zerolog.ErrorLevel
}
func (l *fileLogger) IsFatalEnabled() bool {
	return l.level <= zerolog.FatalLevel
}
func (l *fileLogger) IsPanicEnabled() bool {
	return l.level <= zerolog.PanicLevel
}

func (l *fileLogger) SetLevel(level string) {
	l.level = getLoggingLevel(level)
}
func (l *fileLogger) GetLevel() string {
	return l.level.String()
}

// endregion
