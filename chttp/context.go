package chttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/zsxm/scgo/config"
	"github.com/zsxm/scgo/ctemplate"
	"github.com/zsxm/scgo/data"
	"github.com/zsxm/scgo/funcmap"
	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/session"
	"github.com/zsxm/scgo/tools"
)

const (
	TEMP_SUFFIX = ".html"
)

var temp *template.Template
var allFilesNames []string
var includeFilesNames []string

type Result struct {
	Code    string
	Codemsg string
	Data    interface{}
}

type ResponseData struct {
	Datas []interface{}
	Data  interface{}
	Page  data.Page
	Url   string
}

//当前请求的上下文
type Context struct {
	Response  http.ResponseWriter
	Request   *http.Request
	Params    url.Values
	MultiFile *MultiFile
	Method    string
	Session   session.Interface
	Config    Config
}

func (this *Context) SetHeader(key, val string) {
	this.Response.Header().Set(key, val)
}

func (this *Context) NewResult() Result {
	return Result{
		Code:    "0",
		Codemsg: "ok",
		Data:    "",
	}
}

//获取参数
func (this *Context) GetParams(key string) []string {
	return this.Params[key]
}

//获取参数
func (this *Context) GetParam(key string) string {
	v := this.Params[key]
	if len(v) > 0 {
		return v[0]
	}
	return ""
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
	var err error
	defer func() {
		if err := recover(); err != nil {
			if config.Conf.Debug {
				log.Debug(err, string(debug.Stack()))
			} else {
				log.Info(err)
			}
		}
	}()

	if err != nil {
		log.Info(err)
	}
	dtam := dataToArrayMap(datas)
	dtam.Url = this.Request.URL.String()
	if config.Conf.Debug {
		tmpIncFns := []string{config.Conf.Template.Dir + name + TEMP_SUFFIX}
		tmpIncFns = append(tmpIncFns, includeFilesNames...)
		t := template.New("T").Funcs(funcmap.FuncMap)
		t, err := t.ParseFiles(tmpIncFns...)
		if err != nil {
			log.Error(err)
		}
		li := strings.LastIndex(name, "/")
		if li != -1 {
			name = name[li+1:]
		}
		err = t.ExecuteTemplate(this.Response, name+TEMP_SUFFIX, dtam)
	} else {
		li := strings.LastIndex(name, "/")
		if li != -1 {
			name = name[li+1:]
		}
		err = temp.ExecuteTemplate(this.Response, name+TEMP_SUFFIX, dtam)
	}
	if err != nil {
		log.Error(err)
		http.Error(this.Response, err.Error(), http.StatusInternalServerError)
	}
}

//响应json
func (this *Context) JSON(v interface{}, hasIndent bool) {
	defer func() {
		if err := recover(); err != nil {
			if config.Conf.Debug {
				log.Debug(err, string(debug.Stack()))
			} else {
				log.Info(err)
			}
		}
	}()
	this.SetHeader("Content-Type", "application/text; charset=utf-8")
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
			if config.Conf.Debug {
				log.Debug(err, string(debug.Stack()))
			} else {
				log.Info(err)
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
			if config.Conf.Debug {
				log.Debug(err, string(debug.Stack()))
			} else {
				log.Info(err)
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

func (this *Context) SetCookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
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
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}
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
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			httponly = true
		}
	}
	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}
	this.SetHeader("Set-Cookie", b.String())
}

//write
func (c *Context) Write(v []byte) (int, error) {
	return c.Response.Write(v)
}

//read body
func (c *Context) ReadBody() ([]byte, error) {
	body := c.Request.Body
	defer body.Close()
	return ioutil.ReadAll(body)
}

//Page
func (c *Context) Page() *data.Page {
	page := &data.Page{}
	var pageNo, pageSize int
	if len(c.GetParams("pageNo")) > 0 {
		pageNo = tools.ParseInteger(c.GetParams("pageNo")[0])
	} else {
		pageNo = 1
	}
	if len(c.GetParams("pageSize")) > 0 {
		pageSize = tools.ParseInteger(c.GetParams("pageSize")[0])
	} else {
		pageSize = 10
	}
	page.New(pageNo, pageSize)
	return page
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

func dataToArrayMap(datas interface{}) ResponseData {
	rd := ResponseData{}
	switch datas.(type) {
	case data.EntityBeanInterface:
		bean := datas.(data.EntityBeanInterface)
		if bean != nil {
			rd.Datas = make([]interface{}, 0, 5)
			fieldNames := bean.FieldNames()
			for _, val := range bean.Entitys().Values() {
				mp := make(map[string]string)
				for _, v := range fieldNames.Names() {
					field := val.Field(v)
					if field != nil {
						mp[v] = field.Value()
					}
				}
				rd.Datas = append(rd.Datas, mp)
			}
		}
		break
	case data.EntitysInterface:
		bean := datas.(data.EntitysInterface)
		if bean != nil {
			rd.Datas = make([]interface{}, 0, 5)
			fieldNames := bean.FieldNames()
			for _, val := range bean.Values() {
				mp := make(map[string]string)
				for _, v := range fieldNames.Names() {
					field := val.Field(v)
					if field != nil {
						mp[v] = field.Value()
					}
				}
				rd.Datas = append(rd.Datas, mp)
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
			rd.Datas = append(rd.Datas, mp)
		}
		break
	default:
		rd.Data = datas
	}
	return rd
}

func Init() {
	var err error
	path, err := tools.CurrentDir()
	if err != nil {
		log.Error(err)
	}
	allFilesNames, includeFilesNames = ctemplate.Temps(path+"/"+config.Conf.Template.Dir, config.Conf.Template.Dir, TEMP_SUFFIX)

	if !config.Conf.Debug {
		temp, err = template.New("T").Funcs(funcmap.FuncMap).ParseFiles(allFilesNames...)
		if err != nil {
			log.Error(err)
		}
	}
	log.Info("Init all Templates [ok]")
}
