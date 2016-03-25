# scgo golang 框架

#获取方式：<br/>
##github.com<br/>
*go get github.com/zsxm/scgo<br/>
##clone<br/>
*git clone https://github.com/zsxm/scgo.git<br/>

##简单说明：<br/>
为了方便使用，执行速度，和开发效率，所有反射的功能都由生成的代码替代<br/>
代码生成：github.com/zsxm/scgo/tools/scgen/scgen.exe 代码生成工具包<br/>
生成工具使用方式：<br/>
在实际项目目录结构中<br/>
projectDir<br/>
\--------source<br/>
\----------------module1<br/>
\-----------------------entity<br/>
\-----------------------action<br/>
\-----------------------log<br/>
\-----------------------service<br/>
\----------------module2<br/>
\-----------------------entity<br/>
\-----------------------action<br/>
\-----------------------log<br/>
\-----------------------service<br/>
以下目录说明<br/>
projectDir：项目，source：go源码,module1：模块,entity：结构实体,结构实体目录需要按照一定格式编写<br/>
entity.go示例代码<br/>
`import (`
`	"github.com/zsxm/scgo/data"`
`)`

`//go:generate $GOPATH/src/github.com/zsxm/scgo/tools/scgen/scgen.exe -fileDir=$GOFILE -projectDir=study/app3 -moduleName=chatol` `-goSource=source`<br/>
`//go:@Table value=users`<br/>
`type Message struct {`

	//go:@Column value=u_id
	//go:@Identif
	id data.String

	//go:@Column value=u_name
	name data.String

	//go:@Column value=u_phone
	phone data.String

	//go:@Column value=u_age
	age data.Integer

	tt data.Integer
`}`<br/>
注解说明：<br/>
因为go不支持注解，所以都是以注释形式存在的自定义注解<br/>
固定格式//go:开头<br/>
//go:@Table、//go:@Column、//go:@Identif，分别是结构Message对应的表名,字段对应的列名，和主键字段，目前只支持这些。<br/>
未添加注解的字段是不会映射到数据表中，但是会自动封装数据进去，除了数据表映射功能。
<br/>
注意：需要配置环境变量GOPATH<br/>
-projectDir和-moduleName是需要配置的项目目录和模块名称,其它两个不变<br/>
还需要一个.bat或.sh执行文件放到entity.go同一目录下<br/>
执行文件代码<br/>
`@echo off`<br/>
`echo [INFO] run go generate.`<br/>
`cd %~dp0`<br/>
`call go generate`<br/>
`exit`<br/>
执行该命令后，将会自动生成,action,log,service,和entity_impl.go等封装好的代码。<br/>
自动生成的代码后缀带\_impl的文件内容是一搬不需要改动的，如果改动了，再去执行该命令将会覆盖掉自己写的代码，所以在其它文件中实现。<br/>
<br/>
chttp:<br/>
  *action映射<br/>
  *请求数据绑定和响应数据封装<br/>

chttplib:<br/>
  *http模拟请求发送<br/>
  *文件上传等<br/>

data:<br/>
  *数据库、缓存操作封装<br/>
  *对结构的数据转换<br/>

filter:<br/>
  *过滤器<br/>

logger:<br/>
  *日志输出<br/>

security:<br/>
  *安全<br/>

soap:<br/>
  *webservice<br/>

template:<br/>
  *模版转换<br/>

tools:<br/>
  *框架所有工具包<br/>
  *代码生成包<br/>

##代码生成示例：http://git.oschina.net/snxamdf/study : app3

