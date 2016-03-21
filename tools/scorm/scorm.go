package scorm

import (
	"fmt"
	"go/ast"
	"regexp"
)

//`orm:"column=text,valit=isnotnull"`
var regOrm = regexp.MustCompile(`(orm:)|([a-zA-Z]+=[a-zA-Z]+)`)
var regOrmKeyVal = regexp.MustCompile(`([a-zA-Z]+=[a-zA-Z]+)`)

//`orm:""` 映射
type Orm struct {
}

func (this *Orm) ToOrm(fields []*ast.Field) {
	fmt.Println("---------------------------------------------------------------")
	for _, field := range fields { //遍历字段
		tag := field.Tag

		//fmt.Printf("%+v\n", field.Tag) //取得 `orm:""`
		orms := regOrm.FindAllString(tag.Value, -1)
		if orms[0] == "orm:" {
			fmt.Println(orms, len(orms))
			for _, v := range orms {
				keyVals := regOrmKeyVal.FindAllString(v, -1)
				fmt.Println(keyVals)
			}
		}
	}
}
