package logger_test

import (
	"sync"
	"testing"

	"github.com/zsxm/scgo/logger"
)

func TestLogger(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	logr := logger.NewLogger("test")
	go func() {
		logr.Info("a")
		wg.Done()
	}()
	wg.Wait()
}
