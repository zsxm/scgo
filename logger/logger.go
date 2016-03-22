package logger

import (
	"fmt"
	"sync"
)

func NewLogger(modelName string) *Logger {
	e := &Logger{
		modelName: modelName,
		level:     all,
		msg:       make(chan *msg, 10000),
		logOut:    make(map[string]LoggerInterface),
	}
	e.setLogOut("console", `{"level":6}`)
	e.setLogOut("file", `{"level":6}`)
	go e.startLog()
	return e
}

type msg struct {
	level int
	msg   string
}

type Logger struct {
	modelName string
	lock      sync.Mutex
	level     int
	msg       chan *msg
	logOut    map[string]LoggerInterface
}

func (this *Logger) setLogOut(name string, config string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if logOut, ok := logFuncs[name]; ok {
		lo := logOut()
		err := lo.init(config)
		if err != nil {
			return err
		}
		this.logOut[name] = lo
	} else {
		return fmt.Errorf("logs: unknown logFuncs %q (forgotten Register?)", name)
	}
	return nil
}

func (this *Logger) startLog() {
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
	if this.level > debug {
		this.write(debug, msg...)
	}
}

func (this *Logger) Info(msg ...interface{}) {
	if this.level > info {
		this.write(info, msg...)
	}
}
