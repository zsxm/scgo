package chttp

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/tools"
	"github.com/zsxm/scgo/tools/uuid"
)

type MultiFile struct {
	File       []multipart.File //文件
	DirName    string           //目录名称,模块名称
	FileName   []string         //文件名称
	SumSize    int64            //文件总大小
	Size       []int64          //文件大小
	isUpload   bool             //是否已上传完文件
	FileNameId []string         //上传完文件的名称
}

type Size interface {
	Size() int64
}

func (this *MultiFile) init(fileHeads []*multipart.FileHeader) {
	size := len(fileHeads)
	this.File = make([]multipart.File, 0, size)
	this.FileName = make([]string, 0, size)
	this.Size = make([]int64, 0, size)
	this.FileNameId = make([]string, 0, size)
	for i := 0; i < size; i++ {
		fileHead := fileHeads[i]
		file, err := fileHead.Open()
		if err != nil {
			log.Error(err)
			break
		}
		if s, ok := file.(Size); ok {
			this.SumSize += s.Size()
			this.Size[i] = s.Size()
		}
		this.File[i] = file
		this.FileName[i] = fileHead.Filename
	}
}

//关闭文件
func (this *MultiFile) Close() error {
	if this.isUpload {
		for _, f := range this.File {
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
			src = Conf.UploadPath
		}
		src = src + "/" + this.DirName
		if !tools.Exist(src) {
			err := os.MkdirAll(src, 0666)
			if err != nil {
				log.Error(err)
			}
		}
		for i, file := range this.File {
			index := strings.LastIndex(this.FileName[i], ".")
			if index != -1 {
				ext = this.FileName[i][index:]
			}
			this.FileNameId[i] = uuid.NewV1().String() + ext
			fileName := src + "/" + this.FileNameId[i]
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
	return nil
}
