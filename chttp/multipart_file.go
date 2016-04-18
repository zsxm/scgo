package chttp

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/zsxm/scgo/config"
	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/tools"
	"github.com/zsxm/scgo/tools/uuid"
)

type MultiFile struct {
	file       []multipart.File //文件
	DirName    string           //目录名称,模块名称
	FileName   []string         //文件名称
	SumSize    int64            //文件总大小
	Size       []int64          //文件大小
	isUpload   bool             //是否已上传完文件
	FileNameId []string         //上传完文件的名称
	FileSuffix []string         //文件类型后缀
}

type Size interface {
	Size() int64
}

func (this *MultiFile) Files() []multipart.File {
	return this.file
}

func (this *MultiFile) init(fileHeads []*multipart.FileHeader) {
	size := len(fileHeads)
	this.file = make([]multipart.File, 0, size)
	this.FileName = make([]string, 0, size)
	this.Size = make([]int64, 0, size)
	this.FileNameId = make([]string, 0, size)
	this.FileSuffix = make([]string, 0, size)
	for _, fileHead := range fileHeads {
		file, err := fileHead.Open()
		if err != nil {
			log.Error(err)
			break
		}
		this.file = append(this.file, file)
		this.FileName = append(this.FileName, fileHead.Filename)
		this.FileSuffix = append(this.FileSuffix, fileHead.Filename[strings.Index(fileHead.Filename, "."):])
		if s, ok := file.(Size); ok {
			this.SumSize += s.Size()
			this.Size = append(this.Size, s.Size())
		}
	}

}

//关闭文件
func (this *MultiFile) Close() error {
	if this.isUpload {
		for _, f := range this.file {
			f.Close()
		}
	}
	return nil
}

//上传文件
func (this *MultiFile) Upload(src string) error {
	if this.isUpload {
		var ext string
		if src == "" {
			src = config.Conf.UploadPath
		}
		src = src + "/" + this.DirName
		log.Println("--------------------", src)
		this.DirName = src
		if !tools.Exist(src) {
			err := os.MkdirAll(src, 0666)
			if err != nil {
				log.Error(err)
			}
		}
		for i, file := range this.file {
			fileName := this.FileName[i]
			index := strings.LastIndex(fileName, ".")
			if index != -1 {
				ext = fileName[index:]
			}
			fileNameId := uuid.NewV1().String() + ext
			this.FileNameId = append(this.FileNameId, fileNameId)
			fileName = src + "/" + fileNameId
			defer file.Close()
			buf, err := ioutil.ReadAll(file)
			if err != nil {
				log.Error(err)
				return err
			}
			err = ioutil.WriteFile(fileName, buf, 0666)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	this.isUpload = false
	this.file = nil
	return nil
}
