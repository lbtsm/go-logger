package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// header definition
const (
	FDate         = 1 << iota                  // the date in the local time zone: 2009:01:23
	FTime                                      // the time in the local time zone: 01:23:23
	FMicroseconds                              // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	FLongFile                                  // full file name and line number: /a/b/c/d.go:23
	FShortFile                                 // final file name element and line number: d.go:23. overrides FLongfile
	FStd          = FDate | FTime | FShortFile // std header 2009:01:23-01:23:23-a.go:30
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	level int
	flag  int
	out   io.Writer
}

func NewLogger(level, flag int, filePath, filePrefix, rotateInterval string) (*Logger, error) {
	logger := &Logger{
		level: level,
		flag:  flag,
	}

	err := NewFile(filePath, filePrefix, rotateInterval, logger)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetFlag(flag int) {
	l.flag = flag
}

func (l *Logger) setOutPut(out io.Writer) {
	l.out = out
}

func (l *Logger) Trace(i ...interface{}) {
	if l.level <= TRACE {
		_ = l.outPut(fmt.Sprintln("[trace]", i))
	}
}

func (l *Logger) Debug(i ...interface{}) {
	if l.level <= DEBUG {
		_ = l.outPut(fmt.Sprintln("[debug]", i))
	}
}

func (l *Logger) Info(i ...interface{}) {
	if l.level <= INFO {
		_ = l.outPut(fmt.Sprintln("[info]", i))
	}
}

func (l *Logger) Warn(i ...interface{}) {
	if l.level <= WARN {
		_ = l.outPut(fmt.Sprintln("[warn]", i))
	}
}

func (l *Logger) Error(i ...interface{}) {
	if l.level <= ERROR {
		_ = l.outPut(fmt.Sprintln("[error]", i))
	}
}

// 致命的错误
func (l *Logger) Fatal(i ...interface{}) {
	_ = l.outPut(fmt.Sprintln("[fatal]", i))
	os.Exit(1)
}

func (l *Logger) TraceF(format string, i ...interface{}) {
	if l.level <= TRACE {
		_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	}
}

func (l *Logger) DebugF(format string, i ...interface{}) {
	if l.level <= DEBUG {
		_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	}
}

func (l *Logger) InfoF(format string, i ...interface{}) {
	if l.level <= INFO {
		_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	}
}

func (l *Logger) WarnF(format string, i ...interface{}) {
	if l.level <= WARN {
		_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	}
}

func (l *Logger) ErrorF(format string, i ...interface{}) {
	if l.level <= ERROR {
		_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	}
}

// 致命的错误
func (l *Logger) Fatalf(format string, i ...interface{}) {
	_ = l.outPut(fmt.Sprintf(format+"\n", i...))
	os.Exit(1)
}

func (l *Logger) outPut(s string) error {
	var (
		// todo 这个byte有点low，感觉不高效，每次都需要创建一个，写完之后，找找别人怎么写的,
		data = []byte{}
		t    = time.Now()
	)
	if l.flag&FDate == FDate {
		data = append(data, fmt.Sprintf("%d-%02d-%02d ", t.Year(), t.Month(), t.Day())...)
	}
	if l.flag&FTime == FTime {
		data = append(data, fmt.Sprintf("%02d:%02d:%02d ", t.Hour(), t.Minute(), t.Second())...)
	}
	if l.flag&FMicroseconds == FMicroseconds {
		data = append(data, fmt.Sprintf("%d ", t.Nanosecond()/1000)...)
	}
	/*
		todo 这里也是一个很有意思的点，为什么是2呢
		自己的理解：这里又调用堆栈的，这个函数是trace、debug、info等函数调用的，这是一层，
				  在外面的函数中，调用了trace、debug、info等函数，真又是一层，跳过2层，就是我想要的函数调用信息了
	*/
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		// 获取不到执行的堆栈文件信息，直接输出log
		data = append(data, s...)
		_, err := l.out.Write(data)
		return err
	}
	// 这里不需要过多设计，如果有人设置了longFile又设置shortFile
	if l.flag&FLongFile == FLongFile {
		data = append(data, fmt.Sprintf("%s:%d ", file, line)...)
	}
	if l.flag&FShortFile == FShortFile {
		// file是全名，需要根据/截取，最后的 .go 文件
		short := filepath.Base(file)
		data = append(data, fmt.Sprintf("%s:%d ", short, line)...)
	}

	data = append(data, s...)
	if l.out == nil {
		return nil
	}
	_, err := l.out.Write(data)
	return err
}
