package logger

import (
	"fmt"
	"sync"

	"github.com/zsxm/scgo/tools/cxml"
)

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

type Logger struct {
	modelName string
	lock      sync.Mutex
	level     int
	msg       chan *msg
	logOut    map[string]LoggerInterface
	Logger    loggerXml `xml:"logger"`
}

func NewLogger(modelName string) *Logger {
	e := &Logger{
		modelName: modelName,
		level:     all,
		msg:       make(chan *msg, 10000),
		logOut:    make(map[string]LoggerInterface),
	}
	e.xmlInit()
	if e.Logger.Console {
		e.setOut("console", e.Logger)
	}
	if e.Logger.File {
		e.setOut("file", e.Logger)
	}
	go e.start()
	return e
}

func (this *Logger) xmlInit() error {
	return cxml.Unmarshal(this, xml_path)
}

func (this *Logger) setOut(name string, xml loggerXml) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if logOut, ok := logFuncs[name]; ok {
		lo := logOut()
		err := lo.init(xml)
		if err != nil {
			return err
		}
		this.logOut[name] = lo
	} else {
		return fmt.Errorf("Not found reg func %q (forgotten register ?)", name)
	}
	return nil
}

func (this *Logger) start() {
	for {
		select {
		case m := <-this.msg:
			for _, v := range this.logOut {
				v.write(m.level, m.msg)
			}
		}
	}
}

func (this *Logger) write(level int, v ...interface{}) {
	m := new(msg)
	m.level = level
	m.msg = fmt.Sprint(this.modelName, " ", log_level[level], " ", fmt.Sprint(v...))
	this.msg <- m
}

func (this *Logger) Debug(msg ...interface{}) {
	if this.level >= debug {
		this.write(debug, msg...)
	}
}

func (this *Logger) Info(msg ...interface{}) {
	if this.level >= info {
		this.write(info, msg...)
	}
}

func (this *Logger) Warn(msg ...interface{}) {
	if this.level >= warn {
		this.write(warn, msg...)
	}
}

func (this *Logger) Error(msg ...interface{}) {
	if this.level >= err {
		this.write(err, msg...)
	}
}

func (this *Logger) Fatal(msg ...interface{}) {
	if this.level >= fatal {
		this.write(fatal, msg...)
	}
}
