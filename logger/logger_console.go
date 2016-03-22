package logger

import (
	"encoding/json"
	"log"
	"os"
)

type console struct {
	lg    *log.Logger
	Level int `json:"level"`
}

func newConsole() LoggerInterface {
	cw := &console{
		lg:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Level: all,
	}
	return cw
}

func (this *console) init(config string) error {
	return json.Unmarshal([]byte(config), this)
}

func (this *console) write(level int, msg string) error {
	if this.Level > level {
		log.Println(msg)
	}
	return nil
}

func init() {
	register("console", newConsole)
}
