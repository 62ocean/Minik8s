package hpa

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
	HpaInit() error
	//HpaChangeHandler(eventType object.EventType, rs object.Hpa)
	AddHpa(rs object.Hpa)
	DeleteHpa(rs object.Hpa)
	UpdateHpa(rs object.Hpa)
}

type controller struct {
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

	//创建client对hpa进行增删改操作
	c.client = HTTPClient.CreateHTTPClient(global.ServerHost)

	//初始化当前etcd中的hpa
	err := c.HpaInit()
	if err != nil {
		fmt.Println("[hpa controller] hpa init failed")
		return
	}

	//创建subscribe监听hpa的变化
	c.s, _ = subscriber.NewSubscriber(global.MQHost)
	c.handler = &hpaHandler{
		c: c,
	}
	err = c.s.Subscribe("hpas", subscriber.Handler(c.handler))
	if err != nil {
		fmt.Println("[hpa controller] subscribe hpa failed")
		return
	}

}

func (c *controller) HpaInit() error {
	//得到所有的hpa列表
	response := c.client.Get("/hpas/getAll")
	hpaList := new(map[string]string)
	err := json.Unmarshal([]byte(response), hpaList)
	if err != nil {
		fmt.Println("[hpa controller] unmarshall hpalist failed")
		return err
	}

	// 为当前的所有replicaset都启动一个worker
	for _, value := range *hpaList {
		//fmt.Println(value)
		var hpa object.Hpa
		err := json.Unmarshal([]byte(value), &hpa)
		if err != nil {
			fmt.Println("[hpa controller] unmarshall hpa failed")
			return err
		}
		c.AddHpa(hpa)
	}

	return nil
}

func (c *controller) AddHpa(hpa object.Hpa) {
	log.Print("[hpa controller] create hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	//RSworker := NewWorker(rs)
	//c.workers[rs.Metadata.Uid] = RSworker
	//go RSworker.Start()
}

func (c *controller) DeleteHpa(hpa object.Hpa) {
	log.Print("[hpa controller] delete hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	//RSworker := c.workers[rs.Metadata.Uid]
	//RSworker.Stop()

}
func (c *controller) UpdateHpa(hpa object.Hpa) {
	log.Print("[hpa controller] update hpa: " + hpa.Metadata.Name + "  uid: " + hpa.Metadata.Uid)

	//RSworker := c.workers[rs.Metadata.Uid]
	//RSworker.UpdateReplicaset(rs)
}

func NewController() Controller {
	c := &controller{}
	c.workers = make(map[string]Worker)

	return c
}

// --------------------hpa change handler----------------
type hpaHandler struct {
	c *controller
}

func (h *hpaHandler) Handle(msg []byte) {
	log.Println("[hpa controller] hpa receive msg: " + string(msg))

	var msgObject object.MQMessage
	err := json.Unmarshal(msg, &msgObject)
	if err != nil {
		fmt.Println("[hpa controller] unmarshall msg failed")
		return
	}

	var hpa object.Hpa
	if msgObject.EventType == object.DELETE {
		err = json.Unmarshal([]byte(msgObject.PrevValue), &hpa)
	} else {
		err = json.Unmarshal([]byte(msgObject.Value), &hpa)
	}

	if err != nil {
		fmt.Println("[hpa controller] unmarshall changed hpa failed")
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
