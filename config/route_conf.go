package config

//路由配置
var Conf *Config

type Config struct {
	Static     Mapping   //静态文件 url 映射css、js、image
	Html       Mapping   //html文件 url 映射
	Port       string    //服务商品
	Error404   ErrorPage //404错误页面
	Error500   ErrorPage //500错误页面
	Template   Template  //html模版
	Icon       Icon      //icon
	Debug      bool      //启动模式
	DBCfgPath  string    //数据库配置文件路径
	UploadPath string    //上传文件目录 默认值 upload
	Welcome    string    //默认首页 默认值 /index.html
}

func (this *Config) Init() {
	if this.UploadPath == "" {
		this.UploadPath = "upload"
	}
	if this.Welcome == "" {
		this.Welcome = "/index.html"
	}
}

type Mapping struct {
	Dir    string   //所在目录
	Prefix string   //请求url
	Ext    []string //文件扩展名
}

type ErrorPage struct {
	Url     string `` //页面路径
	Message string `` //返回错误信息
}

type Icon struct {
	Name string "" //icon.ico 名称
}

type Template struct {
	Dir     string  "template" //模版所在根目录
	Include Include //引入的模版配置
}

type Include struct {
	Files []string //包含的模版文件名称
	Dir   string   //所在目录
}
