package main

import (
	"fmt"
	"k8s/pkg/controllers"
	"log"
	"os"
)

func init() {
	//err := os.MkdirAll("../../log/controllerManager/", 755)
	//fmt.Println("aaaaa")
	//if err != nil {
	//	fmt.Println("create dir failed")
	//	return
	//}
	//fmt.Println("aaaaa")
	//
	//TODO log不好用

	//logFile, err := os.OpenFile("log/controllers manage"+time.Now().Format("15_04_05")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logFile, err := os.OpenFile("log/controllers manage.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("[controllers manager]")
}

func main() {
	m := controllers.NewManager()
	m.Start()
}
