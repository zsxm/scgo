package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/zsxm/scgo/tools"
	"github.com/zsxm/scgo/tools/date"
)

type LogFile struct {
	lg       *log.Logger
	mw       *MuxWriter
	FileName string `json:"filename"`
	Level    int    `json:"level"`
	MaxSize  int64  `json:"maxSize"`
	filepath string
}

type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

func (this *MuxWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	return this.fd.Write(b)
}

func NewLogFile() LoggerInterface {
	cw := &LogFile{
		Level: ALL,
	}
	cw.mw = new(MuxWriter)
	cw.lg = log.New(cw.mw, "", log.Ldate|log.Ltime)
	return cw
}

func (this *LogFile) Init(config string) error {
	json.Unmarshal([]byte(config), this)
	if this.MaxSize > 0 {
		maxSize = this.MaxSize
	}
	err := this.createLogFile()
	if err != nil {
		return err
	}
	return this.fileSize()
}

func (this *LogFile) Write(level int, msg string) {
	this.fileSize()
	if this.Level > level {
		this.lg.Println(msg)
	}
}

func (this *LogFile) createLogFile() error {
	if tools.IsNotBlank(this.FileName) {
		fileName = this.FileName + "."
	}
	if !exist(path) {
		os.Mkdir(path, 0660)
	}
	this.filepath = path + fileName + date.FormatNumString(date.FormatYYMD(date.Now())) + ext
	fd, err := os.OpenFile(this.filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	this.mw.fd = fd
	return err
}

func (this *LogFile) fileSize() error {
	this.mw.Lock()
	defer this.mw.Unlock()
	fd := this.mw.fd
	s, _ := fd.Stat()
	if s.Size() > maxSize {
		err := fd.Close()
		if err != nil {
			return err
		}
		fname := path + fileName + date.FormatNumString(date.FormatYYMDHMS(date.Now())) + ext
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

func init() {
	Register("file", NewLogFile)
}
