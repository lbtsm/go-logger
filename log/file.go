package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DefaultRotateInterval = "1h"
	FileFormat            = "%s-%d-%02d-%02d.%02d.%02d.%02d-%d.log"
)

type File struct {
	filePath       string
	filePrefix     string
	logger         *Logger
	outPut         io.Writer
	m              sync.Mutex
	rotateInterval time.Duration
}

func NewFile(filePath, filePrefix, rotateInterval string, logger *Logger) error {
	f := &File{
		filePath:   filePath,
		filePrefix: filePrefix,
		logger:     logger,
	}

	d, err := time.ParseDuration(rotateInterval)
	if err != nil {
		d, _ = time.ParseDuration(DefaultRotateInterval)
	}
	f.rotateInterval = d
	f.createFile()
	logger.setOutPut(f)
	go f.logRotate()

	return nil
}

func (f *File) Write(p []byte) (n int, err error) {
	// todo 每次都写入，会不会很耗性能呢？
	return f.outPut.Write(p)
}

func (f *File) logRotate() {
	for {
		select {
		// 根据设置的轮转时间重新生成文件
		case <-time.After(f.rotateInterval):
			f.createFile()
			// todo 在轮转时间内，文件超大怎么办呢？
		}
	}
}

func (f *File) createFile() {
	now := time.Now()
	fileName := fmt.Sprintf(FileFormat, f.filePrefix, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
	newFilePath := filepath.Join(f.filePath, fileName)
	file, err := os.Create(newFilePath)
	if err != nil {
		log.Fatalln("日志轮转失败，使用旧的IO输出", err)
		return
	}
	// 按照理解来讲，当m.lock()锁住之后，该协程的资源被锁住，等待解锁
	f.m.Lock()
	f.outPut = file
	f.m.Unlock()
}
