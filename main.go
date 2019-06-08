package main

import (
	"time"

	"github.com/lbtsm/go-logger/log"
)

func main() {
	for i := 0; i < 10; i++ {
		go func(flag int) {
			w := log.NewWorker(flag + 1)
			w.Do()
		}(i)
	}

	for {
		time.Sleep(time.Second * 3)
	}
}
