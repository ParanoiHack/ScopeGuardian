package logger

import (
	"golang.org/x/exp/slog"
)

type slogAdapter struct {
	inner *slog.Logger
}

// NewSlogLogger wraps an *slog.Logger and returns it as a Logger.
func NewSlogLogger(inner *slog.Logger) Logger {
	return &slogAdapter{inner: inner}
}

// Info forwards an informational log message to the underlying slog.Logger.
func (a slogAdapter) Info(s string, fields ...Field) {
	a.inner.Info(s, fieldToAttr(fields))
}

// Error forwards an error log message to the underlying slog.Logger.
func (a slogAdapter) Error(s string, fields ...Field) {
	a.inner.Error(s, fieldToAttr(fields))
}

// Panic logs the message as an error via the underlying slog.Logger and then panics.
func (a slogAdapter) Panic(s string, fields ...Field) {
	a.inner.Error(s, fieldToAttr(fields))
	panic(s)
}

// fieldToAttr converts a slice of Field values into slog.Attr values.
func fieldToAttr(fields []Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))

	for _, f := range fields {
		attrs = append(attrs, slog.Any(f.Key, f.Value))
	}

	return attrs
}
