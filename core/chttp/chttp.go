package chttp

import (
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

type action map[string]*curl

type curl struct {
	permissions []string
	mfunc       func(Context)
	method      string
}

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
			if Conf.Debug {
				log.Println(err, string(debug.Stack()))
			} else {
				log.Panicln(err)
			}
			this.Error500(w, r)
		}
	}()
	url := r.URL.String()
	//log.Println("--------- ", url, this.action)
	ix := strings.Index(url, "?")
	if ix > 0 {
		url = url[0:ix]
	}
	if this.isStatic(url) { //*.js、*.css、image等静态文件
		staticRoute.init(w, r)
		return
	} else if this.isHtml(url) { //*.html
		htmlRoute.init(w, r)
		return
	} else if murl, ok := this.action[url]; ok { //请求url判断
		if murl.method == ALL || murl.method == r.Method { //请求方式判断
			//判断Action Url是否设置了url权限
			if murl.permissions != nil && len(murl.permissions) > 0 {

			} else {
				c, err := this.Context(w, r)
				if err != nil {
					log.Println(err)
				}
				murl.mfunc(c) //调用函数
			}
		} else {
			log.Println("未找到这个 URL", url, "URL 请求方式", murl.method, ",当前请求方式", r.Method)
			this.Error404(w, r)
		}
	} else {
		log.Println("未找到这个 URL", url, ",请求方式", r.Method)
		this.Error404(w, r)
	}
}

//500 error
func (*Route) Error500(w http.ResponseWriter, r *http.Request) {
	if Conf.Error500.Url != "" {
		http.Redirect(w, r, Conf.Error500.Url, http.StatusFound)
		return
	} else {
		w.WriteHeader(500)
		w.Write([]byte(Conf.Error500.Message))
	}
}

//404 error
func (*Route) Error404(w http.ResponseWriter, r *http.Request) {
	if Conf.Error404.Url != "" {
		http.Redirect(w, r, Conf.Error404.Url, http.StatusFound)
		return
	} else {
		w.WriteHeader(404)
		w.Write([]byte(Conf.Error404.Message))
	}
}

//判断是否为静态js、css、image等文件请求
func (*Route) isStatic(url string) bool {
	if strings.HasPrefix(url, Conf.Static.Prefix) {
		for _, v := range Conf.Static.Ext {
			if strings.HasSuffix(url, v) {
				return true
			}
		}
	}
	return false
}

//判断是否为静态html文件请求
func (*Route) isHtml(url string) bool {
	if strings.HasPrefix(url, Conf.Html.Prefix) {
		for _, v := range Conf.Html.Ext {
			if strings.HasSuffix(url, v) {
				return true
			}
		}
	}
	return false
}

//当前请求的上下文
func (*Route) Context(w http.ResponseWriter, r *http.Request) (Context, error) {
	values := url.Values{}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return Context{}, err
	}
	if r.Method == GET {
		values = r.Form
	} else {
		values = r.PostForm
	}

	c := Context{
		Response: w,
		Request:  r,
		Params:   values,
	}
	return c, nil
}

//运行服务
func Run() {
	log.Println("HTTP PROT", Conf.Port)
	err := http.ListenAndServe(Conf.Port, route)
	if err != nil {
		log.Println(err)
	}
}
