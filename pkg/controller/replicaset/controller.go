package replicaset

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/controller/replicaset/worker"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"sync"
)

type Controller interface {
	Start(wg *sync.WaitGroup)
	ReplicasetInit() error
	ReplicasetChangeHandler(eventType object.EventType, rs object.ReplicaSet)
	AddReplicaset(rs object.ReplicaSet)
	DeleteReplicaset(rs object.ReplicaSet)
	UpdateReplicaset(rs object.ReplicaSet)
}

type controller struct {
	workers map[string]worker.Worker

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

	//初始化当前etcd中的replicaset
	err := c.ReplicasetInit()
	if err != nil {
		fmt.Println("rs init failed")
		return
	}

	//创建subscribe监听replicaset的变化
	c.s, _ = subscriber.NewSubscriber("amqp://guest:guest@localhost:5672/")
	c.handler = NewReplicasetChangeHandler(c)
	err = c.s.Subscribe("replicasets", subscriber.Handler(c.handler))
	if err != nil {
		fmt.Println("subcribe rs failed")
		return
	}

}

func (c *controller) ReplicasetInit() error {
	//得到所有的rs列表
	response := c.client.Get("/replicasets/getAll")
	rsList := new(map[string]string)
	err := json.Unmarshal([]byte(response), rsList)
	if err != nil {
		fmt.Println("unmarshall podlist failed")
		return err
	}

	// 为当前的所有replicaset都启动一个worker
	for _, value := range *rsList {
		fmt.Println(value)
		var rs object.ReplicaSet
		err := json.Unmarshal([]byte(value), &rs)
		if err != nil {
			fmt.Println("unmarshall rs failed")
			return err
		}
		c.AddReplicaset(rs)
	}

	return nil
}

func (c *controller) ReplicasetChangeHandler(eventType object.EventType, rs object.ReplicaSet) {

	switch eventType {
	case object.CREATE:
		c.AddReplicaset(rs)
	case object.DELETE:
		c.DeleteReplicaset(rs)
	case object.UPDATE:
		c.UpdateReplicaset(rs)
	}
}

func (c *controller) AddReplicaset(rs object.ReplicaSet) {
	fmt.Print("create replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	quit := make(chan int)
	RSworker := worker.NewWorker(rs, quit)
	c.workers[rs.Metadata.Uid] = RSworker
	RSworker.Start()
}

func (c *controller) DeleteReplicaset(rs object.ReplicaSet) {
	fmt.Print("delete replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	RSworker := c.workers[rs.Metadata.Uid]
	RSworker.Stop()

}
func (c *controller) UpdateReplicaset(rs object.ReplicaSet) {
	fmt.Print("update replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	RSworker := c.workers[rs.Metadata.Uid]
	RSworker.UpdateReplicaset(rs)
}

func NewController() Controller {
	c := &controller{}
	c.workers = make(map[string]worker.Worker)

	return c
}
