package filter

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/zsxm/scgo/session"
)

type Filter struct {
	url   string                    //拦截的url
	ffunc func(FilterContext) error //执行的函数
}

var furl []*Filter

type FilterContext struct {
	Response http.ResponseWriter
	Request  *http.Request
	Session  session.Interface
	Params   url.Values
}

//设置头
func (this *FilterContext) SetHeader(key, val string) {
	this.Response.Header().Set(key, val)
}

//获取参数
func (this *FilterContext) GetParam(key string) []string {
	return this.Params[key]
}

func (this *FilterContext) Redirect(url string, status ...int) {
	code := http.StatusFound
	if len(status) == 1 {
		code = status[0]
	}
	http.Redirect(this.Response, this.Request, url, code)
}

//添加过滤器方法
func Add(url string, filterMethod func(FilterContext) error) {
	fu := &Filter{
		url:   url,
		ffunc: filterMethod,
	}
	furl = append(furl, fu)
}

//调用
func Call(curl string, fc FilterContext) error {
	for _, v := range furl {
		if isCall(curl, v.url) {
			err := v.ffunc(fc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isCall(curl, furl string) bool {
	if furl == "/*" {
		return true
	} else {
		x := strings.Index(furl, "*")
		if x != -1 {
			furl = strings.Replace(furl, "*", "", -1)
		}
		x = strings.LastIndex(furl, "/")
		l := len(furl)
		if x == l-1 {
			furl = furl[0 : l-1]
		}
		x = strings.Index(curl, furl)
		if x == 0 {
			return true
		}
	}
	return false
}

func init() {
	furl = make([]*Filter, 0, 1)
}
