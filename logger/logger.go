package logger

import (
	"fmt"
	"log"
	"sync"
)

type loggerFunc func() LoggerInterface

var logFuncs = make(map[string]loggerFunc)

func Register(name string, log loggerFunc) {
	logFuncs[name] = log
}

func NewLogger() *Logger {
	e := &Logger{
		level:  ALL,
		msg:    make(chan *msg, 10000),
		logOut: make(map[string]LoggerInterface),
	}
	go e.startLog()
	return e
}

type msg struct {
	level int
	msg   string
}

type Logger struct {
	lock   sync.Mutex
	level  int
	msg    chan *msg
	logOut map[string]LoggerInterface
}

func (this *Logger) SetLogOut(name string, config string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if logOut, ok := logFuncs[name]; ok {
		lo := logOut()
		err := lo.Init(config)
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
				v.Write(m.level, m.msg)
			}
		}
	}
}

func (this *Logger) write(level int, v ...interface{}) {
	m := new(msg)
	m.level = level
	m.msg = fmt.Sprint(LOG_LEVEL[level], " ", fmt.Sprint(v...))
	this.msg <- m
}

func (this *Logger) Debug(msg ...interface{}) {
	if this.level > DEBUG {
		this.write(DEBUG, msg...)
	}
}

func (this *Logger) Info(msg ...interface{}) {
	if this.level > INFO {
		this.write(INFO, msg...)
	}
}

func init() {
}
