package main

var loggerTemp = `//scgen
package log

import (
	"github.com/zsxm/scgo/logger"
)

var loger *logger.Logger = logger.NewLogger("[{{.GenEntity.ModuleName}}]")

func Debug(msg interface{}) {
	loger.Debug(msg)
}

func Info(msg interface{}) {
	loger.Info(msg)
}
`
