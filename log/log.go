package log

import (
	"github.com/zsxm/scgo/logger"
)

var loger *logger.Log = logger.New("[ZSXM-SCGO]")

func Debug(msg ...interface{}) {
	loger.Debug(msg...)
}

func Info(msg ...interface{}) {
	loger.Info(msg...)
}

func Warn(msg ...interface{}) {
	loger.Warn(msg...)
}

func Error(msg ...interface{}) {
	loger.Error(msg...)
}

func Fatal(msg ...interface{}) {
	loger.Fatal(msg...)
}
