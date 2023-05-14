package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/pkg/api/pod"
	"os"
)
import (
	"k8s/pkg/api/service"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("test pr!")
	service.ServiceConfigTest()
	// 解析pod的yaml配置文件
	dataBytes, err := os.ReadFile("D:\\Homework\\K8s\\repository\\k8s\\pkg\\pod\\podConfigTest.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	var podData pod.Pod
	err2 := yaml.Unmarshal(dataBytes, &podData)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	fmt.Println(podData)

	// 根据配置文件创建容器

}
