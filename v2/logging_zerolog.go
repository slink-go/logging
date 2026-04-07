package logging

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
)

// region - zerolog

func newZerologLogger(id string, opts ...Option) Logger {
	var w io.Writer
	if os.Getenv("GO_ENV") == "dev" || strings.ToLower(os.Getenv("LOGGING_FORMAT")) == "pretty" {
		w = configureConsoleWriter(id)
	} else {
		w = os.Stdout
	}
	ctx := zerolog.New(w).Level(getLoggingLevel(id)).With().Str("logger", id).Timestamp()

	for _, opt := range opts {
		if opt != nil {
			ctx = opt(ctx)
		}
	}

	logger := ctx.Logger()
	result := &zerologLogger{
		lg: &logger,
	}
	return result
}

type zerologLogger struct {
	lg *zerolog.Logger
}

func (l *zerologLogger) Clone(newId string) Logger {
	return GetLogger(newId)
}
func (l *zerologLogger) Trace(message string, args ...interface{}) {
	if l.lg.GetLevel() == zerolog.TraceLevel {
		l.lg.Trace().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Debug(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.DebugLevel {
		l.lg.Debug().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Info(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.InfoLevel {
		l.lg.Info().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Warning(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.WarnLevel {
		l.lg.Warn().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Error(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.ErrorLevel {
		l.lg.Error().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Fatal(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.FatalLevel {
		l.lg.Fatal().Fields(args).Msg(message)
	}
}
func (l *zerologLogger) Panic(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.PanicLevel {
		l.lg.Panic().Fields(args).Msg(message)
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

func (l *zerologLogger) SetLevel(level string) Logger {
	v := l.lg.Level(getLoggingLevel(level))
	l.lg = &v
	return l
}
func (l *zerologLogger) GetLevel() string {
	return l.lg.GetLevel().String()
}

// endregion
