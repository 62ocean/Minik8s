package main

import (
	"fmt"
	"k8s/pkg/apiserver/flannel"
	"k8s/pkg/kubelet"
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
		logFile, err = os.OpenFile("../log/kubelet.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if sysType == "windows" {
		// windows系统
		logFile, err = os.OpenFile("log/kubelet.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("[Kubelet]")
}

func main() {
	flannel.ConfigInit()
	// 创建kubelet对象
	kl, _ := kubelet.NewKubelet("node1")
	defer kubelet.StopKubelet(kl)
	// 创建kubelet对象
	kl.Run()
}
