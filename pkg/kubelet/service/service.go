package service

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"os"
)

func ServiceConfigTest() {
	fmt.Println("测试 service 配置文件读取")
	dataBytes, err := os.ReadFile("pkg/api/service/serviceConfigTest.yaml")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	fmt.Println("文件内容：\n", string(dataBytes))
	service := object.Service{}
	err = yaml.Unmarshal(dataBytes, &service)
	if err != nil {
		fmt.Println("解析yaml文件失败:", err)
		return
	}
	fmt.Printf("解析结果：\n + service -> %+v\n", service)
}

func CreateService(serviceConfig object.Service) {
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	resp := client.Get("/pods/getAll")
	podsBody, _ := json.Marshal(resp)
	var allPods []object.Pod
	json.Unmarshal(podsBody, allPods)
	for _, pod := range allPods {
		if serviceConfig.Spec.Selector.App == pod.Metadata.Labels.App {
			if serviceConfig.Spec.Selector.Env == pod.Metadata.Labels.Env {
				serviceConfig.Spec.Pods = append(serviceConfig.Spec.Pods, pod)
			}
		}
	}

}
