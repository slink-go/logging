// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2024-2026 Mikhail Kantur

package logging

type PyroscopeLoggerImpl struct {
	logger Logger
}

func (l *PyroscopeLoggerImpl) Infof(fmt string, args ...interface{}) {
	l.logger.Info(fmt, args...)
}
func (l *PyroscopeLoggerImpl) Debugf(fmt string, args ...interface{}) {
	l.logger.Debug(fmt, args...)
}
func (l *PyroscopeLoggerImpl) Errorf(fmt string, args ...interface{}) {
	l.logger.Error(fmt, args...)
}

func PyroscopeLogger(logger Logger) *PyroscopeLoggerImpl {
	return &PyroscopeLoggerImpl{
		logger: logger,
	}
}
