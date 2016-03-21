package logger

type LoggerInterface interface {
	Init(config string) error
	Write(level int, msg string)
}
