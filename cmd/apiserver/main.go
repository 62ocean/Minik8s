package main

import (
	"fmt"
	"k8s/pkg/apiserver"
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
		logFile, err = os.OpenFile("../../log/apiserver.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if sysType == "windows" {
		// windows系统
		logFile, err = os.OpenFile("log/apiserver.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("[APIServer]")
}

func main() {

	//etcd.EtcdTest()
	//apiserver.StartServer()
	//kubectl.CmdExec()
	//fmt.Println("hello world")
	//log.Println("test Log!")
	//etcd.EtcdInit(global.EtcdHost)

	server, _ := apiserver.CreateAPIServer()
	server.StartServer()
}
