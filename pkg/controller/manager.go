package controller

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"k8s/pkg/controller/replicaset"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"os"
	"sync"
)

const (
	controllerNum = 1
)

type Manager interface {
	Start()
}

type manager struct {
	replicasetController replicaset.Controller
}

func (m *manager) Start() {

	var wg sync.WaitGroup
	wg.Add(controllerNum)

	fmt.Println("1213123")

	go m.replicasetController.Start(&wg)

	// test: add a replicaset to apiserver
	// --------------------------------------
	dataBytes, err := os.ReadFile("test/ReplicasetConfigTest.yml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	//fmt.Println("aaaaaa")
	var replicasetData object.ReplicaSet
	err2 := yaml.Unmarshal(dataBytes, &replicasetData)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	//err = utils.OutputJson("解析yaml: replicaset", replicasetData)
	if err != nil {
		fmt.Println("解析yaml: replicaset失败")
		return
	}
	id, _ := uuid.NewUUID()
	replicasetData.Metadata.Uid = id.String()
	var rsJson []byte
	rsJson, _ = json.Marshal(replicasetData)
	//fmt.Println("rsJson: \n" + string(rsJson))

	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	client.Post("/replicasets/create", rsJson)
	fmt.Println("add replicaset ok!")
	//--------------------------------------

	// 等待所有协程执行完毕
	wg.Wait()
}

func NewManager() Manager {
	m := &manager{}

	fmt.Print("qqq")

	// 创建各种controller
	m.replicasetController = replicaset.NewController()

	return m
}
