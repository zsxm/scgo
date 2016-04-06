package chttp

import (
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"

	"github.com/zsxm/scgo/config"
	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/session"
)

type action map[string]*curl

//路由
type Route struct {
	action action
}

//路由
var route = &Route{
	action: make(map[string]*curl),
}

//实现ServeHTTP
func (this *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//如果有错误就恢复 并跳转到错误页面
	defer func() {
		if err := recover(); err != nil {
			if config.Conf.Debug {
				log.Debug(err, string(debug.Stack()))
			} else {
				log.Info(err)
			}
			this.Error500(w, r)
		}
	}()
	url := r.URL.String()
	ix := strings.Index(url, "?")
	if ix > 0 {
		url = url[0:ix]
	}
	if url == "/" {
		if this.isHtml(config.Conf.Welcome) {
			htmlRoute.init(w, r)
			return
		}
	}
	if this.isStatic(url) { //*.js、*.css、image等静态文件
		staticRoute.init(w, r)
		return
	} else if this.isHtml(url) { //*.html
		htmlRoute.init(w, r)
		return
	} else if murl, ok := this.action[url]; ok { //请求url判断
		if (murl.method == ALL || murl.method == r.Method) && murl.mtype == MTYPE_HTTP { //请求方式判断 http
			//判断Action Url是否设置了url权限
			if murl.permissions != nil && len(murl.permissions) > 0 {

			} else {
				c, err := this.Context(w, r)
				if err != nil {
					log.Error(err)
				}
				fc := this.FilterContext(c)
				err = Call(url, fc) //调用过滤器
				if err != nil {
					log.Error(err)
					return
				}
				//control配置
				c.SetControlConfig(murl.controlConfig)
				murl.mfunc(c) //调用函数
				multi := c.MultiFile()
				if multi != nil && multi.isUpload {
					var src = config.Conf.UploadPath
					err := multi.Upload(src)
					if err != nil {
						log.Error(err)
					}
					multi.Close()
				}
			}
		} else if murl.method == GET && murl.mtype == MTYPE_WEBSOCKET { //请求方式判断 websocket
			ws := &WebSocketRoute{murl.wshandler}
			ws.ServeHTTP(w, r)
			return
		} else {
			if url != "/favicon.ico" {
				log.Info("未找到 URL ", url, ",请求方式", murl.method, ",当前请求方式", r.Method)
				this.Error404(w, r)
			}
		}
	} else {
		if url != "/favicon.ico" {
			log.Info("未找到 URL ", url, ",请求方式", r.Method)
			this.Error404(w, r)
		}
	}
}

//500 error
func (*Route) Error500(w http.ResponseWriter, r *http.Request) {
	if config.Conf.Error500.Url != "" {
		http.Redirect(w, r, config.Conf.Error500.Url, http.StatusFound)
		return
	} else {
		w.WriteHeader(500)
		w.Write([]byte(config.Conf.Error500.Message))
	}
}

//404 error
func (*Route) Error404(w http.ResponseWriter, r *http.Request) {
	if config.Conf.Error404.Url != "" {
		http.Redirect(w, r, config.Conf.Error404.Url, http.StatusFound)
		return
	} else {
		w.WriteHeader(404)
		w.Write([]byte(config.Conf.Error404.Message))
	}
}

//判断是否为静态js、css、image等文件请求
func (*Route) isStatic(url string) bool {
	if strings.HasPrefix(url, config.Conf.Static.Prefix) {
		for _, v := range config.Conf.Static.Ext {
			if strings.HasSuffix(url, v) {
				return true
			}
		}
	}
	return false
}

//判断是否为静态html文件请求
func (*Route) isHtml(url string) bool {
	if strings.HasPrefix(url, config.Conf.Html.Prefix) {
		for _, v := range config.Conf.Html.Ext {
			if strings.HasSuffix(url, v) {
				return true
			}
		}
	}
	return false
}

//判断是否为静态html文件请求
func (*Route) isWebSocket(r *http.Request) bool {

	return false
}

//过滤器上下文
func (*Route) FilterContext(c Context) FilterContext {
	fc := FilterContext{}
	fc.Context = c
	return fc
}

//当前请求的上下文
func (*Route) Context(w http.ResponseWriter, r *http.Request) (Context, error) {
	values := url.Values{}
	err := r.ParseForm()
	if err != nil {
		log.Info(err)
		return nil, err
	}
	c := &context{}
	if r.Method == GET {
		values = r.Form
	} else {
		var contextType = r.Header.Get("Content-Type")
		var isUpload bool
		if contextType != "" {
			if strings.Contains(contextType, "multipart/form-data") {
				isUpload = true
			}
		}
		if isUpload {
			r.ParseMultipartForm(32 << 20) //32M
			multiFile := &MultiFile{}
			fileHeads := r.MultipartForm.File["file"]
			multiFile.init(fileHeads)
			multiFile.isUpload = true
			c.multiFile = multiFile
			values = r.Form
		} else {
			values = r.PostForm
			if len(values) == 0 {
				values = r.Form
			}
		}
	}
	c.method = r.Method
	c.response = w
	c.request = r
	c.params = values
	c.session = session.New(w, r, session.OptionsConfig)
	return c, nil
}

//运行服务
func Run() {
	config.Conf.Init()
	Init()
	log.Info("HTTP PROT", config.Conf.Port, "[ok]")
	err := http.ListenAndServe(config.Conf.Port, route)
	if err != nil {
		log.Info(err)
	}
}
