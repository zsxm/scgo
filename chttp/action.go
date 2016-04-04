package chttp

type curl struct {
	permissions   []string
	mfunc         func(Context)
	method        string
	controlConfig *ControlConfig
}

type Control struct {
	controlConfig *ControlConfig
}

func NewControl(cc *ControlConfig) *Control {
	return &Control{controlConfig: cc}
}

//设置新的action
func (this *Control) Add(url string, actionMethod func(Context)) *curl {
	if route.action == nil {
		route.action = make(map[string]*curl)
	}
	ml := &curl{
		mfunc:         actionMethod,
		method:        ALL,
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
