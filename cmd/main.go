package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/pkg/controller/replicaset"
	"k8s/utils"
	"os"
	"time"
)

//func init() {
//	logFile, err := os.OpenFile("log/"+time.Now().Format("15_04_05")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//	if err != nil {
//		fmt.Println("open log file failed, err:", err)
//		return
//	}
//	log.SetOutput(logFile)
//	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
//}

func main() {
	//etcd.EtcdTest()
	//apiserver.StartServer()
	//kubectl.CmdExec()
	//fmt.Println("hello world")
	//log.Println("test Log!")

	//解析replicaset.yaml
	dataBytes, err := os.ReadFile("../pkg/replicaset/ReplicasetConfigTest.yml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	var replicasetData replicaset.ReplicaSet
	err2 := yaml.Unmarshal(dataBytes, &replicasetData)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	utils.OutputJson("replicaset", replicasetData)

	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		fmt.Print("每隔3秒执行任务")
	}
	ticker.Stop()

	//fmt.Println(replicasetData)

	// 解析pod的yaml配置文件
	//dataBytes, err := os.ReadFile("../pkg/api/pod/podConfigTest.yaml")
	//if err != nil {
	//	fmt.Println("读取文件失败：", err)
	//	return
	//}
	//var podData pod.Pod
	//err2 := yaml.Unmarshal(dataBytes, &podData)
	//if err2 != nil {
	//	fmt.Println("解析 yaml 文件失败：", err)
	//}
	//fmt.Println(podData)

	//dataBytes, err := os.ReadFile("../pkg/api/pod/podConfigTest.yaml")
	//if err != nil {
	//	fmt.Println("读取文件失败：", err)
	//	return
	//}
	//fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	//var podData pod.Pod
	//err2 := yaml.Unmarshal(dataBytes, &podData)
	//if err2 != nil {
	//	fmt.Println("解析 yaml 文件失败：", err)
	//}
	//fmt.Println(podData)

	//kubectl.CmdExec()

	// 根据配置文件创建容器
	//err = pod.CreatePod(podData)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

}
