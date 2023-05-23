package main

import (
	"fmt"
	"k8s/pkg/kubelet"
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
	log.SetPrefix("[Kubelet]")
}

func main() {
	// 创建kubelet对象
	kl, _ := kubelet.NewKubelet("node1")
	kl.Run()
}
