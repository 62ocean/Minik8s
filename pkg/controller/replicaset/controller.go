package replicaset

import (
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"sync"
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

	//s监听replicaset的变化，handler处理
	s       *subscriber.Subscriber
	handler *RSChangeHandler

	//client通过http进行replicaset的增改删
	client *HTTPClient.Client
}

func (c *controller) Start(wg *sync.WaitGroup) {
	// 该函数退出时将锁-1，确保主进程不会在该协程退出之前退出
	defer wg.Done()

	//创建client对replicaset进行增删改操作
	c.client = HTTPClient.CreateHTTPClient(global.ServerHost)

	//创建subscribe监听replicaset的变化
	c.s, _ = subscriber.NewSubscriber("amqp://guest:guest@localhost:5672/")
	c.handler = NewReplicasetChangeHandler(c)
	err := c.s.Subscribe("replicasets", subscriber.Handler(c.handler))
	if err != nil {
		fmt.Println("subcribe rs failed")
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
	fmt.Print("create replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	//quit := make(chan int)
	//worker := NewWorker(rs, quit)
	//c.workers[rs.Metadata.Uid] = worker
	//go worker.Start()
}

func (c *controller) DeleteReplicaset(rs object.ReplicaSet) {
	fmt.Print("delete replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	worker := c.workers[rs.Metadata.Uid]
	worker.Stop()

}
func (c *controller) UpdateReplicaset(rs object.ReplicaSet) {
	fmt.Print("update replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	worker := c.workers[rs.Metadata.Uid]
	worker.UpdateReplicaset(rs)
}

func NewController() Controller {
	c := &controller{}
	c.workers = make(map[string]Worker)

	return c
}
