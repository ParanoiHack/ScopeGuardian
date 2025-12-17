package logger

import (
	"golang.org/x/exp/slog"
)

type slogAdapter struct {
	inner *slog.Logger
}

func NewSlogLogger(inner *slog.Logger) Logger {
	return &slogAdapter{inner: inner}
}

func (a slogAdapter) Info(s string, fields ...Field) {
	a.inner.Info(s, fieldToAttr(fields))
}

func (a slogAdapter) Error(s string, fields ...Field) {
	a.inner.Error(s, fieldToAttr(fields))
}

func (a slogAdapter) Panic(s string, fields ...Field) {
	a.inner.Error(s, fieldToAttr(fields))
	panic(s)
}

func fieldToAttr(fields []Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))

	for _, f := range fields {
		attrs = append(attrs, slog.Any(f.Key, f.Value))
	}

	return attrs
}
