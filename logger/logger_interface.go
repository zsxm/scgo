package logger

type LoggerInterface interface {
	init(config string) error
	write(level int, msg string) error
}

type loggerFunc func() LoggerInterface

var logFuncs = make(map[string]loggerFunc)

func register(name string, log loggerFunc) {
	logFuncs[name] = log
}
