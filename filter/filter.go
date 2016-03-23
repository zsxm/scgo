package filter

import (
	"net/http"
	"net/url"
	"strings"
)

type Filter struct {
	url   string
	ffunc func(FilterContext) error
}

var furl []*Filter

type FilterContext struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   url.Values
}

func Add(url string, filterMethod func(FilterContext) error) {
	fu := &Filter{
		url:   url,
		ffunc: filterMethod,
	}
	furl = append(furl, fu)
}

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
