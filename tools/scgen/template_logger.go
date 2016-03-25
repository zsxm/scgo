package main

var loggerTemp = `//scgen
package log

import (
	"github.com/zsxm/scgo/logger"
)

var loger *logger.Log = logger.New("[{{if isNotBlank .GenEntity.ModuleName}}{{.GenEntity.ModuleName}}{{else}}{{.GenEntity.ProjectDir}}{{end}}]")

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

func Printf(fmt string, msg ...interface{}) {
	loger.Printf(fmt, msg...)
}

func Println(msg ...interface{}) {
	loger.Println(msg...)
}
`
