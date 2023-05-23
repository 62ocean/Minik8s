package service

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"k8s/pkg/global"
	kube_proxy "k8s/pkg/kube-proxy"
	"k8s/pkg/util/HTTPClient"
	"os"
)

func ServiceConfigTest() {
	fmt.Println("测试 service 配置文件读取")
	dataBytes, err := os.ReadFile("test/serviceConfigTest.yaml")
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

//func ServiceInit() {
//	kube_proxy.KubeProxyInit()
//}

func CreateService(serviceConfig object.Service) {
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	resp := client.Get("/pods/getAll")
	podsBody, _ := json.Marshal(resp)
	// 把selector对应的pod放到service对象里
	var allPods []object.Pod
	json.Unmarshal(podsBody, allPods)
	for _, pod := range allPods {
		if serviceConfig.Spec.Selector.App == pod.Metadata.Labels.App {
			if serviceConfig.Spec.Selector.Env == pod.Metadata.Labels.Env {
				serviceConfig.Spec.Pods = append(serviceConfig.Spec.Pods, pod)
			}
		}
	}
	// 配置iptable
	kube_proxy.RegisterService(serviceConfig)

}

func DeleteePod(service object.Service, pod object.Pod) {
	// 删除当前iptable
	kube_proxy.DeleteService(service)

	// 删除pod
	var index int
	for i, p := range service.Spec.Pods {
		if p.Metadata.Uid == pod.Metadata.Uid {
			index = i
			break
		}
	}
	service.Spec.Pods = append(service.Spec.Pods[:index], service.Spec.Pods[index+1:]...)

	// 重新配置iptable
	kube_proxy.RegisterService(service)
}

func AddPod(service object.Service, pod object.Pod) {
	kube_proxy.DeleteService(service)
	service.Spec.Pods = append(service.Spec.Pods, pod)
	kube_proxy.RegisterService(service)
}
