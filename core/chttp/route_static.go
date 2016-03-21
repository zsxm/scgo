package chttp

import (
	"net/http"
)

//静态路由
var staticRoute *StaticRoute

type StaticRoute struct {
	handle http.Handler
}

//Static路由实现ServeHTTP
func (this *StaticRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handle.ServeHTTP(w, r)
}

//初始化 staticRoute
func (this *StaticRoute) init(w http.ResponseWriter, r *http.Request) {
	if staticRoute == nil {
		staticRoute = &StaticRoute{
			handle: http.StripPrefix(Conf.Static.Prefix, http.FileServer(http.Dir(Conf.Static.Dir))),
		}

	}
	staticRoute.ServeHTTP(w, r)
}
