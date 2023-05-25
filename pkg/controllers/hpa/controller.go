package hpa

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/controllers/replicaset"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"sync"
)

type Controller interface {
	Start(wg *sync.WaitGroup)
	HpaInit() error
	//HpaChangeHandler(eventType object.EventType, rs object.Hpa)
	AddHpa(rs object.Hpa)
	DeleteHpa(rs object.Hpa)
	UpdateHpa(rs object.Hpa)
}

// 使用接口避免包依赖循环
type manager interface {
	GetRSController() replicaset.Controller
}

type controller struct {
	//controller manager的接口， 用于与其他controller的通信
	m manager

	cache Cache

	workers map[string]Worker

	//s监听hpa的变化，handler处理
	s       *subscriber.Subscriber
	handler *hpaHandler

	//client通过http进行hpa的增改删
	client *HTTPClient.Client
}

func (c *controller) Start(wg *sync.WaitGroup) {
	// 该函数退出时将锁-1，确保主进程不会在该协程退出之前退出
	defer wg.Done()

	// cache开始自动同步pod的状态
	go c.cache.SyncLoop()

	//开始监听hpa变化
	err := c.s.Subscribe("hpas", subscriber.Handler(c.handler))
	if err != nil {
		fmt.Println("[hpa controllers] subscribe hpa failed")
		return
	}

}

func (c *controller) HpaInit() error {
	//得到所有的hpa列表
	response := c.client.Get("/hpas/getAll")
	hpaList := new(map[string]string)
	err := json.Unmarshal([]byte(response), hpaList)
	if err != nil {
		fmt.Println("[hpa controllers] unmarshall hpalist failed")
		return err
	}

	// 为当前的所有replicaset都启动一个worker
	for _, value := range *hpaList {
		//fmt.Println(value)
		var hpa object.Hpa
		err := json.Unmarshal([]byte(value), &hpa)
		if err != nil {
			fmt.Println("[hpa controllers] unmarshall hpa failed")
			return err
		}
		c.AddHpa(hpa)
	}

	return nil
}

func (c *controller) AddHpa(hpa object.Hpa) {
	log.Print("[hpa controllers] create hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	RSworkers := c.m.GetRSController().GetAllWorkers()
	targetRSworker, ok := RSworkers[hpa.Spec.ScaleTargetRef.Name]
	if !ok {
		log.Println("[hpa controller] create hpa failed (target rs doesn't exist)")
		return
	}

	HPAworker := NewWorker(hpa, c.cache, targetRSworker, c.client)
	c.workers[hpa.Metadata.Uid] = HPAworker
	go HPAworker.Start()
}

func (c *controller) DeleteHpa(hpa object.Hpa) {
	log.Print("[hpa controllers] delete hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	HPAworker := c.workers[hpa.Metadata.Uid]
	HPAworker.Stop()

}
func (c *controller) UpdateHpa(hpa object.Hpa) {
	log.Print("[hpa controllers] update hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	HPAworker := c.workers[hpa.Metadata.Uid]
	HPAworker.UpdateHpa(hpa)
}

func NewController(manager manager, client *HTTPClient.Client) Controller {
	c := &controller{}
	c.workers = make(map[string]Worker)
	c.m = manager
	c.cache = NewCache(client)
	c.client = client

	//初始化当前etcd中的hpa
	err := c.HpaInit()
	if err != nil {
		fmt.Println("[hpa controllers] hpa init failed")
		return nil
	}

	//创建subscribe监听hpa的变化
	c.s, _ = subscriber.NewSubscriber(global.MQHost)
	c.handler = &hpaHandler{
		c: c,
	}

	return c
}

// --------------------hpa change handler----------------

type hpaHandler struct {
	c *controller
}

func (h *hpaHandler) Handle(msg []byte) {
	log.Println("[hpa controllers] receive hpa change msg")

	var msgObject object.MQMessage
	err := json.Unmarshal(msg, &msgObject)
	if err != nil {
		fmt.Println("[hpa controllers] unmarshall msg failed")
		return
	}

	var hpa object.Hpa
	if msgObject.EventType == object.DELETE {
		err = json.Unmarshal([]byte(msgObject.PrevValue), &hpa)
	} else {
		err = json.Unmarshal([]byte(msgObject.Value), &hpa)
	}

	if err != nil {
		fmt.Println("[hpa controllers] unmarshall changed hpa failed")
		return
	}

	switch msgObject.EventType {
	case object.CREATE:
		h.c.AddHpa(hpa)
	case object.DELETE:
		h.c.DeleteHpa(hpa)
	case object.UPDATE:
		h.c.UpdateHpa(hpa)
	}
}
