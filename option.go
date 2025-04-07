package logging

import "github.com/rs/zerolog"

type Option interface {
	Apply(*zerolog.Logger)
}

type callerOpt struct {
}

func (c *callerOpt) Apply(logger *zerolog.Logger) {
	*logger = logger.With().Caller().Logger()
}

func WithCaller() Option {
	return &callerOpt{}
}
