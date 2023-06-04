package parseYaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"os"
)

func ParseReplicasetYaml(filepath string) object.ReplicaSet {

	dataBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		os.Exit(-1)
	}

	var replicaset object.ReplicaSet
	err2 := yaml.Unmarshal(dataBytes, &replicaset)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err2.Error())
		os.Exit(-1)
	}
	//err = utils.OutputJson("解析yaml: replicaset", replicasetData)
	if err != nil {
		fmt.Println("解析yaml: replicaset失败")
		os.Exit(-1)
	}

	return replicaset
}
