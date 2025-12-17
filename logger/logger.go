package logger

var globalLogger Logger = &nullLogger{}

func SetGlobalLogger(l Logger) {
	globalLogger = l
}

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Info(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
}

func Info(s string, fields ...Field) {
	globalLogger.Info(s, fields...)
}

func Error(s string, fields ...Field) {
	globalLogger.Error(s, fields...)
}

func Panic(s string, fields ...Field) {
	globalLogger.Panic(s, fields...)
}

func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

type nullLogger struct{}

func (n nullLogger) Info(_ string, _ ...Field) {}

func (n nullLogger) Error(_ string, _ ...Field) {}

func (n nullLogger) Panic(s string, _ ...Field) {
	panic(s)
}
