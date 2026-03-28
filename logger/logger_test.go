package logger

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNullLogger_Info(t *testing.T) {
	n := &nullLogger{}
	assert.NotPanics(t, func() {
		n.Info("test message")
		n.Info("with field", Field{Key: "k", Value: "v"})
	})
}

func TestNullLogger_Error(t *testing.T) {
	n := &nullLogger{}
	assert.NotPanics(t, func() {
		n.Error("error message")
		n.Error("with field", Field{Key: "k", Value: 42})
	})
}

func TestNullLogger_Panic(t *testing.T) {
	n := &nullLogger{}
	assert.Panics(t, func() {
		n.Panic("panic message")
	})
}

func TestSetGlobalLogger(t *testing.T) {
	original := globalLogger
	defer func() { globalLogger = original }()

	custom := &nullLogger{}
	SetGlobalLogger(custom)

	assert.Equal(t, custom, globalLogger)
}

func TestPackageLevel_Info(t *testing.T) {
	assert.NotPanics(t, func() {
		Info("info message")
		Info("info with field", Field{Key: "key", Value: "val"})
	})
}

func TestPackageLevel_Error(t *testing.T) {
	assert.NotPanics(t, func() {
		Error("error message")
		Error("error with field", Field{Key: "key", Value: "val"})
	})
}

func TestPackageLevel_Panic(t *testing.T) {
	assert.Panics(t, func() {
		Panic("panic message")
	})
}

func TestErr(t *testing.T) {
	field := Err(assert.AnError)

	assert.Equal(t, "error", field.Key)
	assert.Equal(t, assert.AnError, field.Value)
}

func TestAny(t *testing.T) {
	field := Any("mykey", 123)

	assert.Equal(t, "mykey", field.Key)
	assert.Equal(t, 123, field.Value)
}

func TestSlogAdapter_Info(t *testing.T) {
	inner := slog.New(slog.NewTextHandler(io.Discard, nil))
	adapter := NewSlogLogger(inner)

	assert.NotPanics(t, func() {
		adapter.Info("info message")
		adapter.Info("info with fields", Field{Key: "k", Value: "v"})
	})
}

func TestSlogAdapter_Error(t *testing.T) {
	inner := slog.New(slog.NewTextHandler(io.Discard, nil))
	adapter := NewSlogLogger(inner)

	assert.NotPanics(t, func() {
		adapter.Error("error message")
		adapter.Error("error with fields", Field{Key: "k", Value: "v"})
	})
}

func TestSlogAdapter_Panic(t *testing.T) {
	inner := slog.New(slog.NewTextHandler(io.Discard, nil))
	adapter := NewSlogLogger(inner)

	assert.Panics(t, func() {
		adapter.Panic("panic message")
	})
}

func TestFieldToAttr(t *testing.T) {
	fields := []Field{
		{Key: "a", Value: "b"},
		{Key: "x", Value: 99},
	}

	attrs := fieldToAttr(fields)

	assert.Len(t, attrs, 2)
	assert.Equal(t, "a", attrs[0].Key)
	assert.Equal(t, "x", attrs[1].Key)
}

func TestFieldToAttr_Empty(t *testing.T) {
	attrs := fieldToAttr([]Field{})
	assert.Empty(t, attrs)
}
