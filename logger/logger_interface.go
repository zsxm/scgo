package logger

type LoggerInterface interface {
	init(xml loggerXml) error
	write(level int, msg ...interface{}) error
}
