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

	// 等待所有协程执行完毕
	wg.Wait()
}

func NewManager() Manager {
	m := &manager{}

	// 创建各种controller
	m.replicasetController = replicaset.NewController()

	return m
}
