package controller

import (
	"k8s/pkg/controller/replicaset"
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

	go m.replicasetController.Start(&wg)

	// test: add a replicaset to apiserver
	// --------------------------------------

	//replicasetData := parseYaml.ParseReplicasetYaml("test/ReplicasetConfigTest.yml")
	//id, _ := uuid.NewUUID()
	//replicasetData.Metadata.Uid = id.String()
	//var rsJson []byte
	//rsJson, _ = json.Marshal(replicasetData)
	////fmt.Println("rsJson: \n" + string(rsJson))
	//
	//client := HTTPClient.CreateHTTPClient(global.ServerHost)
	//client.Post("/replicasets/create", rsJson)
	//fmt.Println("add replicaset ok!")
	//--------------------------------------

	// 等待所有协程执行完毕
	wg.Wait()
}

func NewManager() Manager {
	m := &manager{}

	// 创建各种controller
	m.replicasetController = replicaset.NewController()

	return m
}
