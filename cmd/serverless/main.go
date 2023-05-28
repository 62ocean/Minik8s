package main

import (
	"fmt"
	"k8s/pkg/serverless"
	"log"
	"os"
)

func init() {
	var logFile *os.File
	var err error
	logFile, err = os.Create("log/serverless.log")
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	//log.SetPrefix("[serverless]")
}

func main() {
	server, _ := serverless.CreateAPIServer()
	server.StartServer()
}
