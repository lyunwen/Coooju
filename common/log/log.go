package log

import (
	"os"
	"time"
)

//日志方法
func Default(log string) {
	log = "\r\n" + time.Now().Format("2006-01-02 15:04:05") + "|" + log
	f, _ := os.OpenFile("log/default.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	_, _ = f.WriteString(log)
}

func Error(error string) {
	error = "\r\n" + time.Now().Format("2006-01-02 15:04:05") + "|" + error
	f, _ := os.OpenFile("log/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	_, _ = f.WriteString(error)
}
