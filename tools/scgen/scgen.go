package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/snxamdf/scgo/tools/annotation"
	"github.com/snxamdf/scgo/tools/gen"
)

func main() {
	log.Println("main")

	fileDir := flag.String("fileDir", "", "file dir")
	projectDir := flag.String("projectDir", "", "project dir")
	moduleName := flag.String("moduleName", "", "module name")
	goSource := flag.String("goSource", "", "go source dir")
	flag.Parse()
	GO_PATH := os.Getenv("GOPATH")
	if GO_PATH == "" {
		log.Fatalln("gopath is null")
	}
	GO_PATH += "/src"
	log.Println("------go path :", GO_PATH)
	log.Println("------project dir :", *projectDir)
	log.Println("------go source dir:", *goSource)
	log.Println("------module name :", *moduleName)
	log.Println("------file dir :", *fileDir)
	gentity := gen.GenEntity{
		GoPath:      GO_PATH,
		ProjectDir:  *projectDir,
		GoSourceDir: *goSource,
		ModuleName:  *moduleName,
		FileDir:     *fileDir,
	}
	f, err := parser.ParseFile(
		token.NewFileSet(),
		*fileDir,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		log.Println(err)
	}
	//bean转换table
	annot := annotation.Bean{
		GenEntity: &gentity,
	}
	//分析orm
	//orm := scorm.Orm{}
	for _, decl := range f.Decls {
		tdecl, ok := decl.(*ast.GenDecl)
		if !ok || tdecl.Doc == nil {
			continue
		}
		for _, comment := range tdecl.Doc.List { //遍历注释
			annot.TableName(comment.Text) //获得表信息
		}
		annot.BeanName(tdecl.Specs[0].(*ast.TypeSpec).Name.Name)
		sdecl := tdecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
		//获取字段
		fields := sdecl.Fields.List
		annot.ToField(fields) //获得field信息
		//orm.ToOrm(fields)      //orm
	}
	//生成entity impl
	genEntity(strings.ToLower(annot.Name)+"_entity_impl.go", annot)
	//生成Action
	genAction(strings.ToLower(annot.Name)+"_action.go", annot)
	//生成service impl
	genServiceImpl(strings.ToLower(annot.Name)+"_service_impl.go", annot)
	//生成service
	genService(strings.ToLower(annot.Name)+"_service.go", annot)
}
