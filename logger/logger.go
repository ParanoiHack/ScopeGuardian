package logger

var globalLogger Logger = &nullLogger{}

// SetGlobalLogger replaces the package-level logger with l.
// All subsequent calls to Info, Error, and Panic will delegate to l.
func SetGlobalLogger(l Logger) {
	globalLogger = l
}

// Field is a structured key-value pair that can be attached to a log entry.
type Field struct {
	Key   string
	Value any
}

// Logger is the interface that all logging backends must satisfy.
type Logger interface {
	Info(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
}

// Info logs an informational message with optional structured fields.
func Info(s string, fields ...Field) {
	globalLogger.Info(s, fields...)
}

// Error logs an error message with optional structured fields.
func Error(s string, fields ...Field) {
	globalLogger.Error(s, fields...)
}

// Panic logs a message and then panics with the message string.
func Panic(s string, fields ...Field) {
	globalLogger.Panic(s, fields...)
}

// Err is a convenience constructor that creates a Field with key "error".
func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

// Any creates a Field with the given key and arbitrary value.
func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

type nullLogger struct{}

func (n nullLogger) Info(_ string, _ ...Field) {}

func (n nullLogger) Error(_ string, _ ...Field) {}

func (n nullLogger) Panic(s string, _ ...Field) {
	panic(s)
}
