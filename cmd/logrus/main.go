package main

import (
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
	"github.com/sirupsen/logrus"
)

func main() {
	l := my_log.NewLog("app/main")

	// 生成一些日誌消息來測試輸出
	l.WithFields(logrus.Fields{
		"username": "JohnDoe",
		"age":      25,
	}).Infof("New user registered")

	l.WithFields(logrus.Fields{
		"path": "/api/action",
	}).Warnf("API endpoint accessed")

	l.Errorf("%s API endpoint failed", "aaa")
	l.Errorln("aaa", "bbb")
}
