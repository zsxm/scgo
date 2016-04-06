package chttp

import (
	"github.com/zsxm/scgo/websocket"
)

type curl struct {
	permissions   []string
	mfunc         func(Context)
	wshandler     func(*websocket.Conn)
	method        string
	mtype         string
	controlConfig *ControlConfig
}

type Control struct {
	controlConfig *ControlConfig
}

func NewControl(cc *ControlConfig) *Control {
	return &Control{controlConfig: cc}
}

//设置新的action
func (this *Control) AddWS(url string, wshandler func(*websocket.Conn)) *curl {
	if route.action == nil {
		route.action = make(map[string]*curl)
	}
	ml := &curl{
		wshandler:     wshandler,
		method:        GET,
		mtype:         MTYPE_WEBSOCKET,
		controlConfig: this.controlConfig,
	}
	route.action[url] = ml
	return ml
}

//设置新的action
func (this *Control) Add(url string, actionMethod func(Context)) *curl {
	if route.action == nil {
		route.action = make(map[string]*curl)
	}
	ml := &curl{
		mfunc:         actionMethod,
		method:        ALL,
		mtype:         MTYPE_HTTP,
		controlConfig: this.controlConfig,
	}
	route.action[url] = ml
	return ml
}

//设置新的action
func Action(url string, actionMethod func(Context)) *curl {
	if route.action == nil {
		route.action = make(map[string]*curl)
	}
	ml := &curl{
		mfunc:  actionMethod,
		method: ALL,
		mtype:  MTYPE_HTTP,
	}
	route.action[url] = ml
	return ml
}

//设置action url方法为Get方式
func (this *curl) Get() *curl {
	this.method = GET
	return this
}

//设置action url方法为Post方式
func (this *curl) Post() *curl {
	this.method = POST
	return this
}

//设置action url 访问权限
func (this *curl) Permission(permissions ...string) *curl {
	this.permissions = permissions
	return this
}
