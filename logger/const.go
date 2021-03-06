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
	path      = "logs/"         //default logs/
	ext       = ".log"          //default .log
	maxSize   = int64(10 << 20) //default 10M
	fileName  = "log."          //default log.
)
