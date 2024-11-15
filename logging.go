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
	SetLevel(level string)
	GetLevel() string
}

// region - zerolog

type zerologLogger struct {
	lg *zerolog.Logger
}

func (l *zerologLogger) Trace(message string, args ...interface{}) {
	if l.lg.GetLevel() == zerolog.TraceLevel {
		l.lg.Trace().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Debug(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.DebugLevel {
		l.lg.Debug().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Info(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.InfoLevel {
		l.lg.Info().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Warning(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.WarnLevel {
		l.lg.Warn().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Error(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.ErrorLevel {
		l.lg.Error().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Fatal(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.FatalLevel {
		l.lg.Fatal().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *zerologLogger) Panic(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.PanicLevel {
		l.lg.Panic().Msg(fmt.Sprintf(message, args...))
	}
}

func (l *zerologLogger) IsTraceEnabled() bool {
	return l.lg.GetLevel() == zerolog.TraceLevel
}
func (l *zerologLogger) IsDebugEnabled() bool {
	return l.lg.GetLevel() <= zerolog.DebugLevel
}
func (l *zerologLogger) IsInfoEnabled() bool {
	return l.lg.GetLevel() <= zerolog.InfoLevel
}
func (l *zerologLogger) IsWarningEnabled() bool {
	return l.lg.GetLevel() <= zerolog.WarnLevel
}
func (l *zerologLogger) IsErrorEnabled() bool {
	return l.lg.GetLevel() <= zerolog.ErrorLevel
}
func (l *zerologLogger) IsFatalEnabled() bool {
	return l.lg.GetLevel() <= zerolog.FatalLevel
}
func (l *zerologLogger) IsPanicEnabled() bool {
	return l.lg.GetLevel() <= zerolog.PanicLevel
}

func (l *zerologLogger) SetLevel(level string) {
	l.lg.Level(getLoggingLevel(level))
}
func (l *zerologLogger) GetLevel() string {
	return l.lg.GetLevel().String()
}

func init() {
	loggerFactory = loggerFactoryImpl{
		loggers: make(map[string]Logger),
	}
}

var loggerFactory loggerFactoryImpl

type loggerFactoryImpl struct {
	loggers map[string]Logger
	mutex   sync.Mutex
}

func GetLogger(id string) Logger {
	loggerFactory.mutex.Lock()
	defer loggerFactory.mutex.Unlock()
	v, ok := loggerFactory.loggers[id]
	if ok {
		return v
	}
	l := newZerologLogger(id)
	loggerFactory.loggers[id] = l
	return l
}

func newZerologLogger(id string) Logger {
	var w io.Writer
	if os.Getenv("GO_ENV") != "dev" {
		w = os.Stdout
	} else {
		w = configureConsoleWriter(id)
	}
	logger := zerolog.
		New(w).
		Level(getLoggingLevel(id)).
		With().Str("logger", id).Timestamp(). //Caller().
		Logger()
	result := &zerologLogger{
		lg: &logger,
	}
	return result
}

func configureConsoleWriter(id string) io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
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

// endregion
// region - slog

// endregion
// region - common

func getConfiguredLevel(id string) string {
	key := os.Getenv("LOGGING_LEVEL_" + strings.ToUpper(strings.ReplaceAll(id, "-", "_")))
	if strings.TrimSpace(key) == "" {
		key = os.Getenv("LOGGING_LEVEL_ROOT")
	}
	return key
}

// endregion
