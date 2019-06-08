package log

import (
	"fmt"
	"time"
)

var defaultLogger *Logger

func init() {
	defaultLogger, _ = NewLogger(INFO, FDate|FTime|FShortFile, "/temp/log", "worker", "1m")
}

type Worker struct {
	flag int
	l    *Logger
}

func NewWorker(flag int) *Worker {
	return &Worker{
		flag: flag,
		l:    defaultLogger,
	}
}

func (w *Worker) Do() {
	for {
		time.Sleep(time.Second * 1)
		prefix := fmt.Sprintf("flag is %d and now time is ", w.flag)
		w.l.Info(prefix, time.Now())
		w.l.InfoF("flag is %d log_worker do is print info log %v", w.flag, time.Now())
	}
}
