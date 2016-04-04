package chttp

import (
	"net/http"
	"net/url"

	"github.com/zsxm/scgo/data"
	"github.com/zsxm/scgo/session"
)

type Context interface {
	BindData(entity data.EntityInterface)
	Download(file string, filename ...string)
	HTML(name string, datas interface{})
	JSON(v interface{}, hasIndent bool)
	Method() string
	MultiFile() *MultiFile
	NewResult() Result
	SetHeader(key, val string)
	Params(key string) []string
	Param(key string) string
	ParamMaps() url.Values
	Page() *data.Page
	Redirect(url string, status ...int)
	Response() http.ResponseWriter
	Request() *http.Request
	ReadBody() ([]byte, error)
	SetCookie(name string, value string, others ...interface{})
	Session() session.Interface
	SetControlConfig(controlConfig *ControlConfig)
	Write(v []byte) (int, error)
	Xml(data interface{}, hasIndent bool)
}
