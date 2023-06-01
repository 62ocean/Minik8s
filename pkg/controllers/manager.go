package controllers

import (
	_ "encoding/json"
	_ "github.com/google/uuid"
	"k8s/pkg/controllers/hpa"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"

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

	//client向apiserver发送http请求
	client *HTTPClient.Client
}

func (m *Manager) Start() {

	var wg sync.WaitGroup
	wg.Add(controllerNum)

	go m.replicasetController.Start(&wg)
	//注释hpa, 便于调试其他功能
	go m.hpaController.Start(&wg)

	// //test: add a replicaset to apiserver
	// //--------------------------------------

	// replicasetData := parseYaml.ParseReplicasetYaml("test/ReplicasetConfigTest.yaml")
	// id, _ := uuid.NewUUID()
	// replicasetData.Metadata.Uid = id.String()
	// var rsJson []byte
	// rsJson, _ = json.Marshal(replicasetData)
	// //fmt.Println("rsJson: \n" + string(rsJson))

	// client := HTTPClient.CreateHTTPClient(global.ServerHost)
	// client.Post("/replicasets/create", rsJson)
	// fmt.Println("add replicaset ok!")
	// //--------------------------------------

	// //test: add a hpa to apiserver
	// //添加hpa前必须有相应的rs，否则会添加失败
	// //--------------------------------------

	// hpaData := parseYaml.ParseYaml[object.Hpa]("test/hpaConfigTest.yaml")
	// id2, _ := uuid.NewUUID()
	// hpaData.Metadata.Uid = id2.String()
	// var rsJson2 []byte
	// rsJson2, _ = json.Marshal(hpaData)
	// //fmt.Println("rsJson: \n" + string(rsJson))

	// client2 := HTTPClient.CreateHTTPClient(global.ServerHost)
	// client2.Post("/hpas/create", rsJson2)
	// fmt.Println("add hpa ok!")
	// //--------------------------------------

	// 等待所有协程执行完毕
	wg.Wait()
}

func (m *Manager) GetRSController() replicaset.Controller {
	return m.replicasetController
}

func (m *Manager) GetHPAController() hpa.Controller {
	return m.hpaController
}

func NewManager() *Manager {
	m := &Manager{}

	m.client = HTTPClient.CreateHTTPClient(global.ServerHost)

	// 创建各种controller
	m.replicasetController = replicaset.NewController(m.client)
	m.hpaController = hpa.NewController(m, m.client)

	return m
}
