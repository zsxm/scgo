package chttp

import (
	"strings"
)

type Filter struct {
	url   string //拦截的url
	nurl  []string
	ffunc func(FilterContext) error //执行的函数
}

var furl []*Filter

type FilterContext struct {
	Context
}

//添加过滤器方法
//url=拦截的url,filterMethod=调用过滤器函数,nurl=不拦截的url
func Add(url string, filterMethod func(FilterContext) error, nurl ...string) {
	fu := &Filter{
		url:   url,
		nurl:  nurl,
		ffunc: filterMethod,
	}
	furl = append(furl, fu)
}

//调用
func Call(curl string, fc FilterContext) error {
	for _, v := range furl {
		for _, nv := range v.nurl {
			if nv == curl { //不拦截的url直接返回
				return nil
			}
		}
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
