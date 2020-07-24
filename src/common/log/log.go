package log

import (
	"httpPortMapping/src/common/util"
	"io"
	"log"
	"os"
)

var (
	info  *log.Logger
	error *log.Logger
)

func Init() {
	if info != nil {
		return
	}
	logFile, err := os.OpenFile("all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	util.ErrCheck(err)
	info = log.New(io.MultiWriter(os.Stderr, logFile), "info:", log.Ldate|log.Ltime|log.Lshortfile)
	error = log.New(io.MultiWriter(os.Stderr, logFile), "error:", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v interface{}) {
	info.Println(v)
}

func Error(v interface{}) {
	error.Println(v)
}
