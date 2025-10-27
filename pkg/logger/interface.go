package logger

type Field struct {
	Key   string
	Value interface{}
}

type ILogger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}
