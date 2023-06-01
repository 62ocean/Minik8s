package main

import (
	"fmt"
	"k8s/pkg/global"
	"k8s/pkg/scheduler"
	"log"
	"os"
	"runtime"
)

func init() {
	sysType := runtime.GOOS
	var logFile *os.File
	var err error
	if sysType == "linux" || sysType == "darwin" {
		// LINUX系统或者MAC
		logFile, err = os.OpenFile("log/scheduler.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if sysType == "windows" {
		// windows系统
		logFile, err = os.OpenFile("log/scheduler.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
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
