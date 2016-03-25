package logger

import (
	"log"
	"os"
)

type console struct {
	lg    *log.Logger
	level int
}

func newConsole() LoggerInterface {
	cw := &console{
		lg:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		level: all,
	}
	return cw
}

func (this *console) init(xml loggerXml) error {
	this.level = xml.Level
	return nil
}

func (this *console) write(level int, msg ...interface{}) error {
	if this.level >= level {
		log.Println(msg...)
	}
	return nil
}
