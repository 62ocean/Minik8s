package service

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func ServiceConfigTest() {
	fmt.Println("测试 service 配置文件读取")
	dataBytes, err := os.ReadFile("pkg/api/service/serviceConfigTest.yaml")
	if err != nil {
		fmt.Println("读取失败", err)
		return
	}
	service := Service{}
	err = yaml.Unmarshal(dataBytes, &service)
	if err != nil {
		fmt.Println("解析失败")
		return
	}
	fmt.Printf("service -> %+v\n", service)
}
