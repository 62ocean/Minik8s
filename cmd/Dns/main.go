package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"k8s/pkg/Dns"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"os"
)

func main() {
	etcd.EtcdInit(global.EtcdHost)

	dataBytes, err := os.ReadFile("./service1.yaml")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	service1 := object.Service{}
	yaml.Unmarshal(dataBytes, &service1)
	serviceName := service1.Metadata.Name
	key := "/registry/services/" + serviceName
	serviceString, _ := json.Marshal(service1)
	etcd.Put(key, string(serviceString))

	dataBytes, err = os.ReadFile("./service2.yaml")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	service2 := object.Service{}
	yaml.Unmarshal(dataBytes, &service2)
	serviceName = service2.Metadata.Name
	key = "/registry/services/" + serviceName
	serviceString, _ = json.Marshal(service2)
	etcd.Put(key, string(serviceString))

	Dns.CreateDns("./dnsConfigTest.yaml")

}
