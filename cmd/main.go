package main

import (
	"k8s/pkg/apiserver"
	"k8s/pkg/kubectl"
	"fmt"
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
}

func main() {
	//etcd.EtcdTest()
	apiserver.StartServer()
	kubectl.CmdExec()
	fmt.Println("hello world")
	log.Println("test Log!")
}
