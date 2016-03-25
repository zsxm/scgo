package logger

import (
	"log"
	"os"
	"sync"

	"github.com/zsxm/scgo/tools"
	"github.com/zsxm/scgo/tools/date"
)

type logFile struct {
	filepath string
	lg       *log.Logger
	mw       *muxWriter
	fileName string
	level    int
	maxSize  int64
}

type muxWriter struct {
	sync.Mutex
	fd *os.File
}

func (this *muxWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	return this.fd.Write(b)
}

func newLogFile() LoggerInterface {
	cw := &logFile{
		level: all,
	}
	cw.mw = new(muxWriter)
	cw.lg = log.New(cw.mw, "", log.Ldate|log.Ltime)
	return cw
}

func (this *logFile) init(xml loggerXml) error {
	if xml.MaxSize != 0 {
		this.maxSize = xml.MaxSize
	} else {
		this.maxSize = maxSize
	}
	if tools.IsNotBlank(xml.FileName) {
		this.fileName = xml.FileName
	} else {
		this.fileName = fileName
	}
	this.level = xml.Level
	err := this.createLogFile()
	if err != nil {
		return err
	}
	return this.fileSize()
}

func (this *logFile) write(level int, msg ...interface{}) error {
	if this.level >= level {
		this.lg.Println(msg...)
	}
	err := this.fileSize()
	if err != nil {
		return err
	}
	return nil
}

func (this *logFile) createLogFile() error {
	if !exist(path) {
		os.Mkdir(path, 0660)
	}
	this.filepath = path + this.fileName + date.FormatNumString(date.FormatYYMD(date.Now())) + ext
	fd, err := os.OpenFile(this.filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	this.mw.fd = fd
	return err
}

func (this *logFile) fileSize() error {
	this.mw.Lock()
	defer this.mw.Unlock()
	fd := this.mw.fd
	s, err := fd.Stat()
	if err != nil {
		return err
	}
	if s.Size() > this.maxSize {
		err = fd.Sync()
		if err != nil {
			return err
		}
		err = fd.Close()
		if err != nil {
			return err
		}
		fname := path + this.fileName + date.FormatNumString(date.FormatYYMDHMS(date.Now())) + ext
		err = os.Rename(this.filepath, fname)
		if err != nil {
			return err
		}
		this.createLogFile()
	}
	return nil
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
