package main

var entityTemp = `//scgen
package entity

import (
	"bytes"
	"fmt"
	"github.com/snxamdf/scgo/data"
	"strconv"
)

//----------------------{{.Name}}Bean begin--------------------------------------

type {{.Name}}Bean struct {
	bean  *{{.Name}}
	beans *{{.Name}}s
}

func New{{.Name}}Bean() *{{.Name}}Bean {
	e := &{{.Name}}Bean{}
	return e
}

func (this *{{.Name}}Bean) NewEntity() data.EntityInterface {
	return New{{.Name}}()
}

func (this *{{.Name}}Bean) NewEntitys(cap int) data.EntitysInterface {
	return New{{.Name}}s(cap)
}

func (this *{{.Name}}Bean) Entity() data.EntityInterface {
	if this.bean == nil {
		return nil
	}
	return this.bean
}

func (this *{{.Name}}Bean) Entitys() data.EntitysInterface {
	if this.beans == nil {
		return nil
	}
	return this.beans
}

func (this *{{.Name}}Bean) Datas() *{{.Name}}s {
	if this.beans == nil {
		return nil
	}
	return this.beans
}

func (this *{{.Name}}Bean) Table() data.TableInformation {
	return {{lower .Name}}TableInformation
}

func (this *{{.Name}}Bean) FieldNames() data.FieldNames {
	return {{lower .Name}}FieldNames
}

func (this *{{.Name}}Bean) SetEntity(bean data.EntityInterface) {
	this.bean = bean.(*{{.Name}})
}

func (this *{{.Name}}Bean) SetEntitys(beans data.EntitysInterface) {
	this.beans = beans.(*{{.Name}}s)
}

//------------------------------------------------------------

//------------------------------------------------------------
type {{.Name}}s struct {
	datas []data.EntityInterface
	page  *data.Page
}

func New{{.Name}}s(cap int) *{{.Name}}s {
	e := &{{.Name}}s{}
	e.datas = make([]data.EntityInterface, 0, cap)
	return e
}

func (this *{{.Name}}s) SetPage(page *data.Page) {
	this.page = page
}

func (this *{{.Name}}s) Add(e data.EntityInterface) {
	this.datas = append(this.datas, e.(*{{.Name}}))
}

func (this *{{.Name}}s) Values() []data.EntityInterface {
	return this.datas
}

func (this *{{.Name}}s) Len() int {
	return len(this.datas)
}

func (this *{{.Name}}s) Get(i int) *{{.Name}} {
	return this.datas[i].(*{{.Name}})
}

func (this *{{.Name}}s) Table() data.TableInformation {
	return {{lower .Name}}TableInformation
}

func (this *{{.Name}}s) FieldNames() data.FieldNames {
	return {{lower .Name}}FieldNames
}

func (this *{{.Name}}s) JSON() string {
	var wr bytes.Buffer
	wr.WriteString("[")
	for i, v := range this.Values() {
		if i > 0 {
			wr.WriteString(",")
		}
		wr.WriteString(v.JSON())
	}
	wr.WriteString("]")
	return wr.String()
}

//----------------------{{.Name}}Bean end--------------------------------------

//----------------------{{.Name}} begin--------------------------------------
func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}
{{$beanName:=.Name}}{{range $bean:=.Fileld}}
func (this *{{$beanName}}) {{upperFirst $bean.Name}}() *data.{{if equal "int" (fieldType $bean.Type)}}Integer{{else if equal "string" (fieldType $bean.Type)}}String{{end}} {
	return &this.{{$bean.Name}}
}

func (this *{{$beanName}}) Set{{upperFirst $bean.Name}}(value {{if equal "int" (fieldType $bean.Type)}}int{{else}}string{{end}}) {
	this.{{$bean.Name}}.SetValue({{if equal "int" (fieldType $bean.Type)}}strconv.Itoa(value){{else}}value{{end}})
}
{{end}}
func (this *{{.Name}}) SetValue(filedName, value string) {
	this.Field(filedName).SetValue(value)
}

func (this *{{.Name}}) Table() data.TableInformation {
	return {{lower .Name}}TableInformation
}

func (this *{{.Name}}) FieldNames() data.FieldNames {
	return {{lower .Name}}FieldNames
}

func (this *{{.Name}}) Field(filedName string) data.EntityField {
	switch filedName {
	{{range $field:=.Fileld}}case {{if equal (lower $field.Name) (lower $field.Column.Name)}}"{{lower $field.Name}}"{{else}}"{{lower $field.Name}}"{{if isNotBlank $field.Column.Name}}, "{{lower $field.Column.Name}}"{{end}}{{end}}:
		{{if $field.Column.Identif}}this.{{$field.Name}}.SetPrimaryKey(true)
		return this.{{$field.Name}}.StructType(){{else}}return this.{{$field.Name}}.StructType(){{end}}
	{{end}}}
	return nil
}

func (this *{{.Name}}) JSON() string {
	var b bytes.Buffer
	b.WriteString("{"){{range $index,$field := .Fileld}}{{if gt $index 0}}b.WriteString(","){{end}}
	b.WriteString(fmt.Sprintf(` + "`" + `"{{.Name}}":%q` + "`" + `, this.{{.Name}}.Value()))
	{{end}}b.WriteString("}")
	return b.String()
}

//----------------------New{{.Name}} end--------------------------------------

//----------------------init() end--------------------------------------
func init() {
	{{lower .Name}}TableInformation.SetTableName("{{.TabName}}")
	{{lower .Name}}ColumnsArr := []string{
		{{range $field:=.Fileld}}{{$colm:=$field.Column}}{{if isNotBlank $colm.Name}}"{{$colm.Name}}", {{end}}{{end}}
	}
	{{lower .Name}}TableInformation.SetColumns({{lower .Name}}ColumnsArr)
	{{lower .Name}}FieldNamesArr := []string{
		{{range $field:=.Fileld}}"{{$field.Name}}", {{end}}
	}
	{{lower .Name}}FieldNames.SetNames({{lower .Name}}FieldNamesArr)
}

var {{lower .Name}}TableInformation data.TableInformation
var {{lower .Name}}FieldNames data.FieldNames

//----------------------init() end--------------------------------------
`
