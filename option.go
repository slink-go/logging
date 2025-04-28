package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"path/filepath"
)

type Option func(zerolog.Context) zerolog.Context

func WithCaller(skipFrameCount ...int) Option {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return func(ctx zerolog.Context) zerolog.Context {
		if len(skipFrameCount) > 0 {
			ctx = ctx.CallerWithSkipFrameCount(skipFrameCount[0])
		} else {
			ctx = ctx.CallerWithSkipFrameCount(3)
		}
		return ctx
	}
}
func WithHook(hook zerolog.Hook) Option {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.Logger().Hook(hook).With()
	}
}
func WithLevel(level zerolog.Level) Option {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.Logger().Level(level).With()
	}
}
