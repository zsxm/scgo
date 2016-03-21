package main

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/snxamdf/scgo/tools"
	"github.com/snxamdf/scgo/tools/annotation"
	"github.com/snxamdf/scgo/tools/gen"
)

//生成实体实现类 xxxx_entity_impl.go
func genEntity(fileName string, annot annotation.Bean) {
	fout, err := os.Create(fileName)

	defer fout.Close()
	if err != nil {
		log.Println(fileName, err)
		return
	}
	buf := bytes.Buffer{}
	temple := newTmpl(entityTemp)
	temple.Execute(&buf, annot)
	n, err := fout.Write(buf.Bytes())
	log.Println(n, err)
}

//生成action类 xxxx_action.go
func genAction(fileName string, annot annotation.Bean) {
	gentity := annot.GenEntity
	genPath := gentity.GoPath + "/" + gentity.ProjectDir + "/" + gentity.GoSourceDir + "/" + gentity.ModuleName + "/"
	log.Println(genPath, exist(genPath))
	fileDir := genPath + "/" + gen.GEN_ACTION
	log.Println(fileDir, exist(fileDir))
	filePath := fileDir + "/" + fileName
	log.Println(filePath, exist(filePath))
	if !exist(fileDir) {
		err := os.MkdirAll(fileDir, 0777)
		log.Println("MkdirAll ", fileDir, err)
	}
	if !exist(filePath) {
		fout, err := os.Create(filePath)
		defer fout.Close()
		if err != nil {
			log.Println(filePath, err)
			return
		}
		buf := bytes.Buffer{}
		temple := newTmpl(actionTemp)
		temple.Execute(&buf, annot)
		n, err := fout.Write(buf.Bytes())
		log.Println(n, err)
	}
}

//生成service类 xxxx_service_impl.go
func genServiceImpl(fileName string, annot annotation.Bean) {
	gentity := annot.GenEntity
	genPath := gentity.GoPath + "/" + gentity.ProjectDir + "/" + gentity.GoSourceDir + "/" + gentity.ModuleName
	log.Println(genPath, exist(genPath))
	fileDir := genPath + "/" + gen.GEN_SERVICE
	log.Println(fileDir, exist(fileDir))
	filePath := fileDir + "/" + fileName
	log.Println(filePath, exist(filePath))
	if !exist(fileDir) {
		err := os.MkdirAll(fileDir, 0777)
		log.Println("MkdirAll ", fileDir, err)
	}
	fout, err := os.Create(filePath)
	defer fout.Close()
	if err != nil {
		log.Println(filePath, err)
		return
	}
	buf := bytes.Buffer{}
	temple := newTmpl(serviceTempImpl)
	temple.Execute(&buf, annot)
	n, err := fout.Write(buf.Bytes())
	log.Println(n, err)
}

//生成service类 xxxx_service.go
func genService(fileName string, annot annotation.Bean) {
	gentity := annot.GenEntity
	genPath := gentity.GoPath + "/" + gentity.ProjectDir + "/" + gentity.GoSourceDir + "/" + gentity.ModuleName
	log.Println(genPath, exist(genPath))
	fileDir := genPath + "/" + gen.GEN_SERVICE
	log.Println(fileDir, exist(fileDir))
	filePath := fileDir + "/" + fileName
	log.Println(filePath, exist(filePath))
	if !exist(fileDir) {
		err := os.MkdirAll(fileDir, 0777)
		log.Println("MkdirAll ", fileDir, err)
	}
	if !exist(filePath) {
		fout, err := os.Create(filePath)
		defer fout.Close()
		if err != nil {
			log.Println(filePath, err)
			return
		}
		buf := bytes.Buffer{}
		temple := newTmpl(serviceTemp)
		temple.Execute(&buf, annot)
		n, err := fout.Write(buf.Bytes())
		log.Println(n, err)
	}
}

//创建一个新模版
func newTmpl(s string) *template.Template {
	return template.Must(template.New("T").Funcs(tools.FuncMap).Parse(s))
}

//判断文件是否存在
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
