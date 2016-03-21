package annotation

import (
	"go/ast"
	"log"
	"regexp"
	"strings"

	"github.com/snxamdf/scgo/tools/gen"
)

//实体映射表
type BeanToTable struct {
	Bean *Bean
}

type Bean struct {
	GenEntity *gen.GenEntity
	TabName   string
	Name      string    //entity名称
	Fileld    *[]Fileld //bean信息
}

type Fileld struct {
	Name   string //字段名
	Type   string //字段类型
	Column Column //列名
}

//表->列
type Column struct {
	Name    string //列名
	Identif bool   //是否唯一标识
}

var regComment, _ = regexp.Compile(`go:@Table|go:@Column|go:@Identif|(value=[\S]+)`)

func (this *Bean) ToField(fields []*ast.Field) {
	log.Println("------ToField------")
	filelds := make([]Fileld, len(fields))
	for i, field := range fields { //遍历字段
		//log.Println(field.Type.(*ast.SelectorExpr).Sel.Name, field.Names)
		filelds[i].Name = field.Names[0].Name                     //字段名称
		filelds[i].Type = field.Type.(*ast.SelectorExpr).Sel.Name //字段类型
		if field.Doc == nil {
			continue
		}
		column := Column{}
		comments := field.Doc.List         //当前字段注解
		for _, comment := range comments { //遍历当前字段的注解
			annot := regComment.FindAllString(comment.Text, -1)
			array := strings.Split(annot[0], ":@")
			if array[0] == "go" {
				switch array[1] {
				case "Identif":
					column.Identif = true
					break
				case "Column":
					column.Name = strings.Split(annot[1], "=")[1]
					break
				}
			}
		}
		filelds[i].Column = column
	}
	this.Fileld = &filelds
	log.Printf("%+v\n", filelds)
}

func (this *Bean) BeanName(comment string) {
	this.Name = comment
}

func (this *Bean) TableName(comment string) {
	annot := regComment.FindAllString(comment, -1)
	if len(annot) > 0 {
		array := strings.Split(annot[0], ":@")
		if array[0] == "go" {
			switch array[1] {
			case "Table":
				this.TabName = strings.Split(annot[1], "=")[1]
				break
			}
		}
	}
}
