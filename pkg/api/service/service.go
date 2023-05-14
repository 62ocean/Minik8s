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
		fmt.Println("读取文件失败:", err)
		return
	}
	fmt.Println("文件内容：\n", string(dataBytes))
	service := Service{}
	err = yaml.Unmarshal(dataBytes, &service)
	if err != nil {
		fmt.Println("解析yaml文件失败:", err)
		return
	}
	fmt.Printf("解析结果：\n + service -> %+v\n", service)
}
