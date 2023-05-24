package controllers

import (
	_ "encoding/json"
	_ "github.com/google/uuid"
	"k8s/pkg/controllers/hpa"
	//"k8s/object"
	"k8s/pkg/controllers/replicaset"
	_ "k8s/pkg/global"
	_ "k8s/pkg/util/HTTPClient"
	_ "k8s/pkg/util/parseYaml"
	_ "log"
	"sync"
)

const (
	controllerNum = 2
)

type Manager struct {
	replicasetController replicaset.Controller
	hpaController        hpa.Controller
}

func (m *Manager) Start() {

	var wg sync.WaitGroup
	wg.Add(controllerNum)

	go m.replicasetController.Start(&wg)
	go m.hpaController.Start(&wg)

	//test: add a replicaset to apiserver
	//--------------------------------------

	//replicasetData := parseYaml.ParseReplicasetYaml("test/ReplicasetConfigTest.yaml")
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

	//test: add a hpa to apiserver
	//--------------------------------------

	//hpaData := parseYaml.ParseYaml[object.Hpa]("test/hpaConfigTest.yaml")
	//id, _ := uuid.NewUUID()
	//hpaData.Metadata.Uid = id.String()
	//var rsJson []byte
	//rsJson, _ = json.Marshal(hpaData)
	////fmt.Println("rsJson: \n" + string(rsJson))
	//
	//client := HTTPClient.CreateHTTPClient(global.ServerHost)
	//client.Post("/hpas/create", rsJson)
	//fmt.Println("add hpa ok!")
	//--------------------------------------

	// 等待所有协程执行完毕
	wg.Wait()
}

func (m *Manager) GetrsController() replicaset.Controller {
	return m.replicasetController
}

func NewManager() *Manager {
	m := &Manager{}

	// 创建各种controller
	m.replicasetController = replicaset.NewController()
	m.hpaController = hpa.NewController(m)

	return m
}
