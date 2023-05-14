package main

import (
	"fmt"
	"k8s/pkg/controller"
	"log"
	"os"
	"time"
)

func init() {
	//err := os.MkdirAll("../../log/controllerManager/", 755)
	//fmt.Println("aaaaa")
	//if err != nil {
	//	fmt.Println("create dir failed")
	//	return
	//}
	//fmt.Println("aaaaa")
	logFile, err := os.OpenFile("log/controllerManager/"+time.Now().Format("15_04_05")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("[controller manager]")
}

func main() {
	log.Println("aaaaaa")
	m := controller.NewManager()
	m.Start()
}
