package logging

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger interface {
	Clone(newId string) Logger
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	Panic(message string, args ...interface{})
	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarningEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
	IsPanicEnabled() bool
	SetLevel(level string) Logger
	GetLevel() string
}

func init() {
	loggerFactory = loggerFactoryImpl{
		consoleLoggers: make(map[string]Logger),
		fileLoggers:    make(map[string]Logger),
	}
}

// region - common

var loggerFactory loggerFactoryImpl

type loggerFactoryImpl struct {
	consoleLoggers map[string]Logger
	fileLoggers    map[string]Logger
	mutex          sync.RWMutex
}

func GetNoOpLogger() Logger {
	return &noOpLogger{}
}
func GetLogger(id string, opts ...Option) Logger {
	loggerFactory.mutex.RLock()
	v, ok := loggerFactory.consoleLoggers[id]
	loggerFactory.mutex.RUnlock()
	if ok {
		return v
	}
	l := newZerologLogger(id, opts...)
	loggerFactory.mutex.Lock()
	loggerFactory.consoleLoggers[id] = l
	loggerFactory.mutex.Unlock()
	return l
}
func DeleteLogger(id string) {
	loggerFactory.mutex.Lock()
	delete(loggerFactory.consoleLoggers, id)
	loggerFactory.mutex.Unlock()
}
func GetFileLogger(file *os.File, id string) Logger {
	loggerFactory.mutex.Lock()
	defer loggerFactory.mutex.Unlock()
	v, ok := loggerFactory.fileLoggers[id]
	if ok {
		return v
	}
	l := newCommonFileLogger(file, id)
	loggerFactory.consoleLoggers[id] = l
	return l
}
func DeleteFileLogger(id string) {
	loggerFactory.mutex.Lock()
	defer loggerFactory.mutex.Unlock()
	v, ok := loggerFactory.fileLoggers[id]
	if ok {
		v.(*fileLogger).Close()
		delete(loggerFactory.fileLoggers, id)
	}
}

// endregion

// region - common

func getConfiguredLevel(id string) string {
	key := os.Getenv("LOGGING_LEVEL_" + strings.ToUpper(strings.ReplaceAll(id, "-", "_")))
	if strings.TrimSpace(key) == "" {
		key = os.Getenv("LOGGING_LEVEL_ROOT")
	}
	return key
}
func configureConsoleWriter(id string) io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		//FormatLevel: func(i interface{}) string {
		//	return strings.ToUpper(fmt.Sprintf("[%5s]", i))
		//},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("[%s] %s", id, i)
		},
		//FormatCaller: func(i interface{}) string {
		//	return filepath.Base(fmt.Sprintf("%s", i))
		//},
		FieldsExclude: []string{
			"logger",
		},
		PartsExclude: []string{
			//zerolog.TimestampFieldName,
			zerolog.CallerFieldName,
			"logger",
		},
	}
}
func getLoggingLevel(id string) zerolog.Level {
	/*
		LevelOffValue = "off" // special value to turn logs off
		LevelTraceValue = "trace"
		LevelDebugValue = "debug"
		LevelInfoValue = "info"
		LevelWarnValue = "warn"
		LevelErrorValue = "error"
		LevelFatalValue = "fatal"
		LevelPanicValue = "panic"
	*/
	key := getConfiguredLevel(id)
	lvl, err := stringToLevel(key)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	return lvl
}
func stringToLevel(key string) (zerolog.Level, error) {
	if strings.ToUpper(key) == "OFF" {
		return zerolog.NoLevel, nil
	}
	if key == "" {
		return zerolog.NoLevel, errors.New("empty config")
	}
	return zerolog.ParseLevel(key)
}
func logLevelAbbr(level zerolog.Level) string {
	switch level {
	case zerolog.TraceLevel:
		return "TRC"
	case zerolog.DebugLevel:
		return "DBG"
	case zerolog.InfoLevel:
		return "INF"
	case zerolog.WarnLevel:
		return "WRN"
	case zerolog.ErrorLevel:
		return "ERR"
	case zerolog.FatalLevel:
		return "FAT"
	case zerolog.PanicLevel:
		return "PNC"
	case zerolog.NoLevel:
		return "NON"
	default:
		return "NON"
	}
}

// endregion
