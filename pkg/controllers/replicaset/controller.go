package replicaset

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"sync"
)

type Controller interface {
	Start(wg *sync.WaitGroup)
	ReplicasetInit() error
	ReplicasetChangeHandler(eventType object.EventType, rs object.ReplicaSet)
	AddReplicaset(rs object.ReplicaSet)
	DeleteReplicaset(rs object.ReplicaSet)
	UpdateReplicaset(rs object.ReplicaSet)

	GetAllWorkers() map[string]Worker
}

type controller struct {
	//每一个rs对应一个worke，存储在workers中
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

	// 开始监听rs变化
	err := c.s.Subscribe("replicasets", subscriber.Handler(c.handler))
	if err != nil {
		fmt.Println("[rs controller] subcribe rs failed")
		return
	}

}

func (c *controller) ReplicasetInit() error {
	//得到所有的rs列表
	response := c.client.Get("/replicasets/getAll")
	rsList := new(map[string]string)
	err := json.Unmarshal([]byte(response), rsList)
	if err != nil {
		fmt.Println("[rs controller] unmarshall rslist failed")
		return err
	}

	// 为当前的所有replicaset都启动一个worker
	for _, value := range *rsList {
		//fmt.Println(value)
		var rs object.ReplicaSet
		err := json.Unmarshal([]byte(value), &rs)
		if err != nil {
			fmt.Println("[rs controller] unmarshall rs failed")
			return err
		}
		c.AddReplicaset(rs)
	}

	return nil
}

func (c *controller) ReplicasetChangeHandler(eventType object.EventType, rs object.ReplicaSet) {

	//fmt.Print(rs)

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
	log.Println("[rs controller] create replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)
	_, ok := c.workers[rs.Metadata.Name]
	if ok {
		log.Println("[rs controller] create replicaset failed! (replicaset name already exists in the same namespace)")
		return
	}

	RSworker := NewWorker(rs, c.client)
	c.workers[rs.Metadata.Name] = RSworker
	go RSworker.Start()
}

func (c *controller) DeleteReplicaset(rs object.ReplicaSet) {
	log.Println("[rs controller] delete replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	RSworker := c.workers[rs.Metadata.Name]
	RSworker.Stop()

}
func (c *controller) UpdateReplicaset(rs object.ReplicaSet) {
	log.Println("[rs controller] update replicaset: " + rs.Metadata.Name + "  uid: " + rs.Metadata.Uid)

	RSworker := c.workers[rs.Metadata.Name]
	RSworker.UpdateReplicaset(rs)
}

func (c *controller) GetAllWorkers() map[string]Worker {
	return c.workers
}

func NewController(client *HTTPClient.Client) Controller {
	c := &controller{}
	c.workers = make(map[string]Worker)
	c.client = client

	//初始化当前etcd中的replicaset
	err := c.ReplicasetInit()
	if err != nil {
		fmt.Println("[rs controller] rs init failed")
		return nil
	}

	//创建subscribe监听replicaset的变化
	c.s, _ = subscriber.NewSubscriber(global.MQHost)
	c.handler = NewReplicasetChangeHandler(c)

	return c
}
