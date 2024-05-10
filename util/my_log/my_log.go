package my_log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf("%s [%s] %s %s:%d\n", timestamp, entry.Level, entry.Message, entry.Caller.File, entry.Caller.Line)), nil
}

func NewLog(filePath string) *logrus.Logger {
	log := logrus.New()
	fileName := fmt.Sprintf("./log/%s.log", filePath)
	logFile := &lumberjack.Logger{
		Filename:   fileName, // 日誌文件路徑
		MaxSize:    2 * 1024, // 檔案最大大小（MB），2GB 需要設定為 2 * 1024
		MaxBackups: 2,        // 保留舊文件的最大個數
		MaxAge:     28,       // 保留舊文件的最大天數
		Compress:   true,     // 是否壓縮/彙縮舊文件
	}

	// 設置輸出為 MultiWriter
	mw := io.MultiWriter(os.Stdout, logFile)

	// 選擇日誌格式（這裡使用 JSON 格式，也可以使用 TextFormatter 等其他格式）
	log.SetFormatter(new(CustomFormatter))
	log.SetOutput(mw)
	//log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)
	return log
}
