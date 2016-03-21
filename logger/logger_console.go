package logger

import (
	"encoding/json"
	"log"
	"os"
)

type Console struct {
	lg    *log.Logger
	Level int `json:"level"`
}

func NewConsole() LoggerInterface {
	cw := &Console{
		lg:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Level: ALL,
	}
	return cw
}

func (this *Console) Init(config string) error {
	return json.Unmarshal([]byte(config), this)
}

func (this *Console) Write(level int, msg string) {
	if this.Level > level {
		log.Println(msg)
	}
}

func init() {
	Register("console", NewConsole)
}
