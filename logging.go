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

type DxLogger interface {
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

type dxLogger struct {
	lg *zerolog.Logger
}

func (l *dxLogger) Trace(message string, args ...interface{}) {
	if l.lg.GetLevel() == zerolog.TraceLevel {
		l.lg.Trace().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Debug(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.DebugLevel {
		l.lg.Debug().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Info(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.InfoLevel {
		l.lg.Info().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Warning(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.WarnLevel {
		l.lg.Warn().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Error(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.ErrorLevel {
		l.lg.Error().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Fatal(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.FatalLevel {
		l.lg.Fatal().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *dxLogger) Panic(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.PanicLevel {
		l.lg.Panic().Msg(fmt.Sprintf(message, args...))
	}
}

func (l *dxLogger) IsTraceEnabled() bool {
	return l.lg.GetLevel() == zerolog.TraceLevel
}
func (l *dxLogger) IsDebugEnabled() bool {
	return l.lg.GetLevel() <= zerolog.DebugLevel
}
func (l *dxLogger) IsInfoEnabled() bool {
	return l.lg.GetLevel() <= zerolog.InfoLevel
}
func (l *dxLogger) IsWarningEnabled() bool {
	return l.lg.GetLevel() <= zerolog.WarnLevel
}
func (l *dxLogger) IsErrorEnabled() bool {
	return l.lg.GetLevel() <= zerolog.ErrorLevel
}
func (l *dxLogger) IsFatalEnabled() bool {
	return l.lg.GetLevel() <= zerolog.FatalLevel
}
func (l *dxLogger) IsPanicEnabled() bool {
	return l.lg.GetLevel() <= zerolog.PanicLevel
}

func (l *dxLogger) SetLevel(level string) {
	l.lg.Level(getLoggingLevel(level))
}
func (l *dxLogger) GetLevel() string {
	return l.lg.GetLevel().String()
}

func init() {
	loggerFactory = loggerFactoryImpl{
		loggers: make(map[string]DxLogger),
	}
}

var loggerFactory loggerFactoryImpl
var mutex sync.RWMutex

type loggerFactoryImpl struct {
	loggers map[string]DxLogger
}

func GetLogger(id string) DxLogger {
	mutex.Lock()
	defer mutex.Unlock()
	v, ok := loggerFactory.loggers[id]
	if ok {
		return v
	}
	l := newLogger(id)
	loggerFactory.loggers[id] = l
	return l
}

func newLogger(id string) DxLogger {
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
	result := &dxLogger{
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
func getConfiguredLevel(id string) string {
	key := os.Getenv("LOGGING_LEVEL_" + strings.ToUpper(strings.ReplaceAll(id, "-", "_")))
	if strings.TrimSpace(key) == "" {
		key = os.Getenv("LOGGING_LEVEL_ROOT")
	}
	return key
}
