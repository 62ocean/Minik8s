package main

import (
	"fmt"
	"k8s/pkg/global"
	"k8s/pkg/scheduler"
	"log"
	"os"
	"time"
)

func init() {
	logFile, err := os.OpenFile("log/"+time.Now().Format("15_04_05")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("[Scheduler]")
}

func main() {
	// 创建scheduler对象
	scheduler, _ := scheduler.NewScheduler(global.ROUND_ROBIN)
	scheduler.Run()
}
