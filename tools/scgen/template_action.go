package main

var actionTemp = `//scgen
package action

import (
	"github.com/zsxm/scgo/chttp"
	"{{.GenEntity.ProjectDir}}/{{.GenEntity.GoSourceDir}}{{if isNotBlank .GenEntity.ModuleName}}/{{.GenEntity.ModuleName}}{{end}}/entity"
)

func init() {
	chttp.Action("{{if isNotBlank .GenEntity.ModuleName}}/{{.GenEntity.ModuleName}}{{end}}/index", index).Get()
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

var controlConf = chttp.ControlConfigConfig()
var control = chttp.NewControl()

func init() {
	controlConf.SetProject("{{.GenEntity.ProjectDir}}")
	controlConf.SetModule("{{.GenEntity.ModuleName}}")
	controlConf.SetTitle("{{.Title}}")
	controlConf.SetComment("{{.Comment}}")
	control.Init(controlConf)
}
`
