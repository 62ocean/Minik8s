package replicaset

import (
	"fmt"
	"k8s/object"
	"k8s/pkg/util/msgQueue/subscriber"
	"sync"
)

const (
	RS_CREATE = iota
	RS_DELETE
	RS_UPDATE
)

type Controller interface {
	Start(wg *sync.WaitGroup)
	//ReplicasetChangeHandler()
	AddReplicaset(rs object.ReplicaSet)
	DeleteReplicaset(rs object.ReplicaSet)
	UpdateReplicaset(rs object.ReplicaSet)
}

type controller struct {
	workers map[string]Worker

	//监听replicaset的变化并处理
	s       *subscriber.Subscriber
	handler *RSChangeHandler
}

func (c *controller) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	c.s, _ = subscriber.NewSubscriber("amqp://guest:guest@localhost:5672/")
	c.handler = NewReplicasetChangeHandler(c)

	err := c.s.Subscribe("replicasets", subscriber.Handler(c.handler))
	if err != nil {
		return
	}

}

//func (c *controller) ReplicasetChangeHandler(msg []byte) {
//
//	// 设msg.rs为发生变化的replicaset, msg.type为发生变化的类型
//	var msg_type int
//	var msg_rs object.ReplicaSet
//
//	switch msg_type {
//	case RS_CREATE:
//		c.AddReplicaset(msg_rs)
//	case RS_DELETE:
//		c.DeleteReplicaset(msg_rs)
//	case RS_UPDATE:
//		c.UpdateReplicaset(msg_rs)
//	}
//}

func (c *controller) AddReplicaset(rs object.ReplicaSet) {
	fmt.Print("create new replicaset: %s %s", rs.Metadata.Name, rs.Metadata.Uid)

	quit := make(chan int)
	worker := NewWorker(rs, quit)
	c.workers[rs.Metadata.Uid] = worker
	go worker.Start()
}

func (c *controller) DeleteReplicaset(rs object.ReplicaSet) {
	fmt.Print("delete replicaset: %s %s", rs.Metadata.Name, rs.Metadata.Uid)

	worker := c.workers[rs.Metadata.Uid]
	worker.Stop()

}
func (c *controller) UpdateReplicaset(rs object.ReplicaSet) {
	fmt.Print("update replicaset: %s %s", rs.Metadata.Name, rs.Metadata.Uid)

	worker := c.workers[rs.Metadata.Uid]
	worker.UpdateReplicaset(rs)
}

func NewController() Controller {
	c := &controller{}
	c.workers = make(map[string]Worker)

	return c
}
