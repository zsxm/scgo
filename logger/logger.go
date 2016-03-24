package logger

import (
	"fmt"
	"log"
	"sync"

	"github.com/zsxm/scgo/tools/cxml"
)

var logger *Logger

type Log struct {
	modelName string
	logger    *Logger
}

type Logger struct {
	lock   sync.Mutex
	level  int
	msg    chan *msg
	logOut map[string]LoggerInterface
	Logger loggerXml `xml:"logger"`
}

type msg struct {
	level int
	msg   string
}

type loggerXml struct {
	Console  bool   `xml:"console"`
	File     bool   `xml:"file"`
	FileName string `xml:"fileName"`
	Level    int    `xml:"level"`
	MaxSize  int64  `xml:"maxSize"`
}

func New(modelName string) *Log {
	e := &Log{
		modelName: modelName,
		logger:    logger,
	}
	return e
}

func (this *Logger) xmlInit() error {
	return cxml.Unmarshal(this, xml_path)
}

func (this *Logger) start() {
	for {
		select {
		case m := <-this.msg:
			for _, v := range this.logOut {
				err := v.write(m.level, m.msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (this *Log) write(level int, v ...interface{}) {
	m := new(msg)
	m.level = level
	m.msg = fmt.Sprint(this.modelName, " ", log_level[level], " ", fmt.Sprint(v...))
	this.logger.msg <- m
}

func (this *Log) Debug(msg ...interface{}) {
	if this.logger.level >= debug {
		this.write(debug, msg...)
	}
}

func (this *Log) Info(msg ...interface{}) {
	if this.logger.level >= info {
		this.write(info, msg...)
	}
}

func (this *Log) Warn(msg ...interface{}) {
	if this.logger.level >= warn {
		this.write(warn, msg...)
	}
}

func (this *Log) Error(msg ...interface{}) {
	if this.logger.level >= err {
		this.write(err, msg...)
	}
}

func (this *Log) Fatal(msg ...interface{}) {
	if this.logger.level >= fatal {
		this.write(fatal, msg...)
	}
}

func init() {
	loggerImpl := make(map[string]LoggerInterface)
	e := &Logger{
		level: all,
		msg:   make(chan *msg, 10240),
	}
	e.xmlInit()
	if e.Logger.Console {
		cons := newConsole()
		cons.init(e.Logger)
		loggerImpl["console"] = cons
	}
	if e.Logger.File {
		fl := newLogFile()
		fl.init(e.Logger)
		loggerImpl["file"] = fl
	}
	e.logOut = loggerImpl
	go e.start()
	logger = e
}
