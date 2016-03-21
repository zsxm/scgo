package main

var actionTemp = `//scgen
package action

import (
	"github.com/zsxm/scgo/chttp"
	"{{.GenEntity.ProjectDir}}/{{.GenEntity.GoSourceDir}}/{{.GenEntity.ModuleName}}/entity"
)

func init() {
	chttp.Action("/{{.GenEntity.ModuleName}}/index", index).Get()
}

//gen
func index(c chttp.Context) {
	e := entity.New{{.Name}}()
	c.BindData(e)
	c.JSON(e.JSON(), true)
}
`
