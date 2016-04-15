package chttp

import (
	"net/http"

	"github.com/zsxm/scgo/config"
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
			handle: http.StripPrefix(config.Conf.Static.Prefix, http.FileServer(http.Dir(config.Conf.Static.Dir))),
		}
	}
	staticRoute.ServeHTTP(w, r)
}

//上传文件
var uploadRoute *UploadRoute

//上传文件
type UploadRoute struct {
	handle http.Handler
}

//Upload路由实现ServeHTTP
func (this *UploadRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handle.ServeHTTP(w, r)
}

//初始化 文件服务uploadRoute
func (this *UploadRoute) init(w http.ResponseWriter, r *http.Request) {
	if uploadRoute == nil {
		path := config.Conf.UploadPath
		if path[0:1] != "/" {
			path = "/" + path
		}
		uploadRoute = &UploadRoute{
			handle: http.StripPrefix(path, http.FileServer(http.Dir(config.Conf.UploadPath))),
		}
	}
	uploadRoute.ServeHTTP(w, r)
}
