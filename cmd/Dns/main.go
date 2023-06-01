package main

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/Dns"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	etcd.EtcdInit(global.EtcdHost)
	client := HTTPClient.CreateHTTPClient(global.ServerHost)

	dataBytes, err := os.ReadFile("./service1.yaml")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	service1 := object.Service{}
	yaml.Unmarshal(dataBytes, &service1)
	getMsg, _ := json.Marshal(service1)
	client.Post("/services/create", getMsg)

	dataBytes, err = os.ReadFile("./service2.yaml")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	service2 := object.Service{}
	yaml.Unmarshal(dataBytes, &service2)
	getMsg, _ = json.Marshal(service2)
	client.Post("/services/create", getMsg)

	Dns.CreateDns()

}
