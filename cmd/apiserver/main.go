package main

import (
	"fmt"
	"k8s/pkg/apiserver"
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
	log.SetPrefix("[APIServer]")
}

func main() {
	//etcd.EtcdTest()
	server := apiserver.CreateAPIServer()
	server.StartServer()
}
