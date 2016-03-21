package chttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snxamdf/scgo/data"
	"github.com/snxamdf/scgo/tools"
)

const (
	TEMP_SUFFIX = ".html"
)

var temp = template.Template{}

//当前请求的上下文
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   url.Values
}

func (this *Context) SetHeader(key, val string) {
	this.Response.Header().Set(key, val)
}

//获取参数
func (this *Context) GetParam(key string) []string {
	return this.Params[key]
}

//绑定实体数据
func (this *Context) BindData(entity data.EntityInterface) {
	for k, v := range this.Params {
		var b bytes.Buffer
		for i, v := range v {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(v)
		}
		field := entity.Field(k)
		if field != nil {
			field.SetValue(b.String())
		}
	}
}

//跳转html模版页面
func (this *Context) HTML(name string, datas interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if Conf.Debug {
				log.Println(err, string(debug.Stack()))
			} else {
				log.Panicln(err)
			}
		}
	}()
	t, err := template.ParseFiles(Conf.Template.Dir + name + TEMP_SUFFIX)
	if err != nil {
		log.Println(err)
	}
	//this.SetHeader("Content-Type", "application/html; charset=utf-8")
	dtam := dataToArrayMap(datas)
	err = t.Execute(this.Response, dtam)
	if err != nil {
		http.Error(this.Response, err.Error(), http.StatusInternalServerError)
	}
}

type ResponseData struct {
	Data []interface{}
	Page data.Page
}

func dataToArrayMap(datas interface{}) ResponseData {
	rd := ResponseData{}
	switch datas.(type) {
	case data.EntityBeanInterface:
		bean := datas.(data.EntityBeanInterface)
		if bean != nil {
			rd.Data = make([]interface{}, 0, 5)
			fieldNames := bean.FieldNames()
			for _, val := range bean.Entitys().Values() {
				mp := make(map[string]string)
				for _, v := range fieldNames.Names() {
					field := val.Field(v)
					if field != nil {
						mp[v] = field.Value()
					}
				}
				rd.Data = append(rd.Data, mp)
			}
		}
		break
	case data.EntitysInterface:
		bean := datas.(data.EntitysInterface)
		if bean != nil {
			rd.Data = make([]interface{}, 0, 5)
			fieldNames := bean.FieldNames()
			for _, val := range bean.Values() {
				mp := make(map[string]string)
				for _, v := range fieldNames.Names() {
					field := val.Field(v)
					if field != nil {
						mp[v] = field.Value()
					}
				}
				rd.Data = append(rd.Data, mp)
			}
		}
		break
	case data.EntityInterface:
		bean := datas.(data.EntityInterface)
		if bean != nil {
			fieldNames := bean.FieldNames()
			mp := make(map[string]string)
			for _, v := range fieldNames.Names() {
				field := bean.Field(v)
				mp[v] = field.Value()
			}
			rd.Data = append(rd.Data, mp)
		}
		break
	}
	return rd
}

//响应json
func (this *Context) JSON(v interface{}, hasIndent bool) {
	defer func() {
		if err := recover(); err != nil {
			if Conf.Debug {
				log.Println(err, string(debug.Stack()))
			} else {
				log.Panicln(err)
			}
		}
	}()
	this.SetHeader("Content-Type", "application/json; charset=utf-8")
	switch v.(type) {
	case string:
		_, err := this.Response.Write([]byte(v.(string)))
		//err := json.NewEncoder(this.Response).Encode(v)
		if err != nil {
			http.Error(this.Response, err.Error(), http.StatusInternalServerError)
		}
		break
	default:
		var content []byte
		var err error
		if hasIndent {
			content, err = json.MarshalIndent(v, "", "  ")
		} else {
			content, err = json.Marshal(v)
		}
		if err != nil {
			http.Error(this.Response, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = this.Response.Write(content)
		if err != nil {
			http.Error(this.Response, err.Error(), http.StatusInternalServerError)
		}
	}
}

//xml
func (this *Context) Xml(data interface{}, hasIndent bool) {
	defer func() {
		if err := recover(); err != nil {
			if Conf.Debug {
				log.Println(err, string(debug.Stack()))
			} else {
				log.Panicln(err)
			}
		}
	}()
	var content []byte
	var err error
	if hasIndent {
		content, err = xml.MarshalIndent(data, "", "  ")
	} else {
		content, err = xml.Marshal(data)
	}
	this.SetHeader("Content-Type", "application/xml; charset=utf-8")
	this.SetHeader("Content-Length", strconv.Itoa(len(content)))
	_, err = this.Response.Write(content)
	if err != nil {
		http.Error(this.Response, err.Error(), http.StatusInternalServerError)
	}
}

//下载
func (this *Context) Download(file string, filename ...string) {
	defer func() {
		if err := recover(); err != nil {
			if Conf.Debug {
				log.Println(err, string(debug.Stack()))
			} else {
				log.Panicln(err)
			}
		}
	}()
	this.SetHeader("Content-Description", "File Transfer")
	this.SetHeader("Content-Type", "application/octet-stream")
	if len(filename) > 0 && filename[0] != "" {
		this.SetHeader("Content-Disposition", "attachment; filename="+filename[0])
	} else {
		this.SetHeader("Content-Disposition", "attachment; filename="+filepath.Base(file))
	}
	this.SetHeader("Content-Transfer-Encoding", "binary")
	this.SetHeader("Expires", "0")
	this.SetHeader("Cache-Control", "must-revalidate")
	this.SetHeader("Pragma", "public")
	http.ServeFile(this.Response, this.Request, file)
}

//重定向
func (this *Context) Redirect(url string, status ...int) {
	code := http.StatusFound
	if len(status) == 1 {
		code = status[0]
	}
	http.Redirect(this.Response, this.Request, url, code)
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}
func (this *Context) Cookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	//fix cookie not work in IE
	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			if v > 0 {
				fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int64:
			if v > 0 {
				fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int32:
			if v > 0 {
				fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}
	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set
	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}
	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}
	// default empty
	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}
	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}
	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}
	this.SetHeader("Set-Cookie", b.String())
}

//Page
func (c *Context) Page() *data.Page {
	page := &data.Page{}
	var pageNo, pageSize int
	if len(c.GetParam("pageNo")) > 0 {
		pageNo = tools.ParseInteger(c.GetParam("pageNo")[0])
	} else {
		pageNo = 1
	}
	if len(c.GetParam("pageSize")) > 0 {
		pageSize = tools.ParseInteger(c.GetParam("pageSize")[0])
	} else {
		pageSize = 10
	}
	page.New(pageNo, pageSize)
	return page
}
