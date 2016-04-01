package chttp

import (
	"github.com/zsxm/scgo/data"
)

type Ctx interface {
	SetHeader(key, val string)
	NewResult() Result
	GetParams(key string) []string
	GetParam(key string) string
	BindData(entity data.EntityInterface)
	HTML(name string, datas interface{})
	JSON(v interface{}, hasIndent bool)
	Xml(data interface{}, hasIndent bool)
	Download(file string, filename ...string)
	Redirect(url string, status ...int)
	SetCookie(name string, value string, others ...interface{})
	Write(v []byte) (int, error)
	ReadBody() ([]byte, error)
	Page() *data.Page
}
