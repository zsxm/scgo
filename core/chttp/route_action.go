package chttp

import (
	"net/http"
)

//静态路由
var actionRoute *ActionRoute

type ActionRoute struct {
	handle http.Handler
}

//Static路由实现ServeHTTP
func (this *ActionRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handle.ServeHTTP(w, r)
}
