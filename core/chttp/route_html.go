package chttp

import (
	"net/http"
)

var htmlRoute *HtmlRoute

type HtmlRoute struct {
	handle http.Handler
}

//Html路由实现ServeHTTP
func (this *HtmlRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handle.ServeHTTP(w, r)
}

//初始化 htmlRoute
func (this *HtmlRoute) init(w http.ResponseWriter, r *http.Request) {
	if htmlRoute == nil {
		htmlRoute = &HtmlRoute{
			handle: http.StripPrefix(Conf.Html.Prefix, http.FileServer(http.Dir(Conf.Html.Dir))),
		}
	}
	htmlRoute.ServeHTTP(w, r)
}
