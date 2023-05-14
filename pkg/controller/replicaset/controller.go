package replicaset

import (
	"fmt"
	"sync"
)

const (
	RS_CREATE = iota
	RS_DELETE
	RS_UPDATE
)

type Controller interface {
	Start(wg *sync.WaitGroup)
	ReplicasetChangeHandler()
	AddReplicaset(rs ReplicaSet)
	DeleteReplicaset(rs ReplicaSet)
	UpdateReplicaset(rs ReplicaSet)
}

type controller struct {
	workers map[string]Worker
}

func (c *controller) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// watch(topic_replicaset, ReplicasetChangeHandler)

}

func (c *controller) ReplicasetChangeHandler() {

	// 设msg.rs为发生变化的replicaset, msg.type为发生变化的类型
	var msg_type int
	var msg_rs ReplicaSet

	switch msg_type {
	case RS_CREATE:
		c.AddReplicaset(msg_rs)
	case RS_DELETE:
		c.DeleteReplicaset(msg_rs)
	case RS_UPDATE:
		c.UpdateReplicaset(msg_rs)
	}
}

func (c *controller) AddReplicaset(rs ReplicaSet) {
	fmt.Print("create new replicaset: %s %s", rs.Metadata.Name, rs.Metadata.uuid)

	quit := make(chan int)
	worker := NewWorker(rs, quit)
	c.workers[rs.Metadata.uuid] = worker
	go worker.Start()
}

func (c *controller) DeleteReplicaset(rs ReplicaSet) {
	fmt.Print("delete replicaset: %s %s", rs.Metadata.Name, rs.Metadata.uuid)

	worker := c.workers[rs.Metadata.uuid]
	worker.Stop()

}
func (c *controller) UpdateReplicaset(rs ReplicaSet) {
	fmt.Print("update replicaset: %s %s", rs.Metadata.Name, rs.Metadata.uuid)

	worker := c.workers[rs.Metadata.uuid]
	worker.UpdateReplicaset(rs)
}

func NewController() Controller {
	c := &controller{}

	return c
}
