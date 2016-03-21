package logger_test

import (
	"sync"
	"testing"

	"time"

	"github.com/zsxm/scgo/logger"
)

func TestLogger(t *testing.T) {
	wg := sync.WaitGroup{}
	var logger = logger.NewLogger()
	logger.SetLogOut("console")
	wg.Add(1)
	go func() {
		time.Sleep(time.Second / 2)
		logger.Debug("fffffffffffff")
		logger.Info("adfasdfasddf")
		wg.Done()
	}()
	wg.Wait()
}
