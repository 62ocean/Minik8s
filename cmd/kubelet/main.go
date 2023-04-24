package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	pod2 "k8s/pkg/api/pod"
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
	log.SetPrefix("[Pod]")
}

func main() {
	// 解析pod的yaml配置文件
	dataBytes, err := os.ReadFile("pkg/api/kubelet/pod/podConfigTest.yaml")
	var podData pod2.Pod
	err2 := yaml.Unmarshal(dataBytes, &podData)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	fmt.Println(podData)

	// 根据配置文件创建容器
	err = pod2.CreatePod(podData)
	if err != nil {
		fmt.Println(err.Error())
	}
}
