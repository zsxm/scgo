package log

import (
	"github.com/zsxm/scgo/logger"
)

var loger *logger.Logger = logger.NewLogger("[ZSXM-SCGO]")

func Debug(msg ...interface{}) {
	loger.Debug(msg...)
}

func Info(msg ...interface{}) {
	loger.Info(msg...)
}
