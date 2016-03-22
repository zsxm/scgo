package logger

const (
	off = iota
	fatal
	err
	warn
	info
	debug
	all
)
const (
	xml_path = "conf/logger.xml"
)

var (
	log_level = []string{"OFF", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "ALL"}
	path      = "logs/"
	ext       = ".log"
	maxSize   = int64(10485760)
	fileName  = "log."
)
