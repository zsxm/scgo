package main

var actionTemp = `//scgen
package action

import (
	"github.com/zsxm/scgo/chttp"
	"{{.GenEntity.ProjectDir}}/{{.GenEntity.GoSourceDir}}{{if isNotBlank .GenEntity.ModuleName}}/{{.GenEntity.ModuleName}}{{end}}/entity"
)

func init() {
	control.Add("{{if isNotBlank .GenEntity.ModuleName}}/{{.GenEntity.ModuleName}}{{end}}/index", index).Get()
}

//gen
func index(c chttp.Context) {
	e := entity.New{{.Name}}()
	c.BindData(e)
	c.JSON(e.JSON(), true)
}
`

var controlTemp = `//scgen
package action

import (
	"github.com/zsxm/scgo/chttp"
)

var controlConf *chttp.ControlConfig = &chttp.ControlConfig{
	Project: "{{.GenEntity.ProjectDir}}",
	Module:  "{{.GenEntity.ModuleName}}",
	Title:   "{{.Title}}",
	Comment: "{{.Comment}}",
}
var control = chttp.NewControl(controlConf)
`
