package logFile

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type LogFile interface {
	Trace() *log.Logger
	Info() *log.Logger
	Error() *log.Logger
}

type logFile struct {
	trace *log.Logger
	info  *log.Logger
	error *log.Logger
}

// NewLogFile log配置
func NewLogFile(dirPath string, fileName string) (logfile LogFile) {
	newPath := filepath.Join("./log", dirPath)
	_ = os.MkdirAll(newPath, os.ModePerm)
	newPath = filepath.Join(newPath, fileName)
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}

	trace := log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	info := log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	e := log.New(io.MultiWriter(file, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logfile = &logFile{
		trace: trace,
		info:  info,
		error: e,
	}
	return
}

func (ft *logFile) Trace() *log.Logger {
	return ft.trace
}

func (ft *logFile) Info() *log.Logger {
	return ft.info
}

func (ft *logFile) Error() *log.Logger {
	return ft.error
}
