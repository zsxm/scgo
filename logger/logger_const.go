package logger

const (
	OFF = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
	ALL
)

var (
	LOG_LEVEL = []string{"OFF", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "ALL"}
	path      = "logs/"
	ext       = ".log"
	maxSize   = int64(10485760)
	fileName  = "log."
)
