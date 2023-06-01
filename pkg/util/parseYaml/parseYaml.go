package parseYaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

//type ObjectType interface {
//	object.Pod | object.Hpa | object.ReplicaSet
//}

func ParseYaml[T any](filepath string) T {

	dataBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		os.Exit(-1)
	}

	var retObject T
	err2 := yaml.Unmarshal(dataBytes, &retObject)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		os.Exit(-1)
	}

	return retObject
}
