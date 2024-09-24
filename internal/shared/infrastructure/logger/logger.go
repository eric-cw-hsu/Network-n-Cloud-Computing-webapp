package logger

type Logger interface {
	Info(arg ...interface{})
	Warn(arg ...interface{})
	Error(arg ...interface{})
	Debug(args ...interface{})
}
