package cron

import (
	"github.com/zsxm/scgo/log"
)

type task struct {
	values map[string]func()
}

var t *task = &task{values: make(map[string]func())}

func Add(name string, method func()) {
	t.values[name] = method
}

func Init() {
	for k, v := range t.values {
		v()
		log.Println("Task Cron Start", k)
		delete(t.values, k)
	}
	t = nil
}
