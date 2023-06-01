package scheduler

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"strconv"
	"sync"
)

type Scheduler struct {
	client      *HTTPClient.Client
	policy      global.Policy
	subscriber  *subscriber.Subscriber
	podHandler  podHandler
	podQueue    string
	nodeHandler nodeHandler
	nodeQueue   string
	nodeList    []object.NodeStorage
}

// NewScheduler Scheduler对象的构造函数
func NewScheduler(p global.Policy) (*Scheduler, error) {
	// 使用HTTP，构建node对象传递到APIServer处
	client := HTTPClient.CreateHTTPClient(global.ServerHost)

	// 建立消息监听队列
	sub, _ := subscriber.NewSubscriber(global.MQHost)

	// 创建scheduler监听队列
	sched := Scheduler{
		client:     client,
		policy:     p,
		subscriber: sub,
		podQueue:   "pods_sched",
		nodeQueue:  "nodes",
	}

	podH := podHandler{
		sched: &sched,
	}
	nodeH := nodeHandler{
		sched: &sched,
	}

	sched.podHandler = podH
	sched.nodeHandler = nodeH
	return &sched, nil
}

// Run Scheduler运行的入口函数
func (s *Scheduler) Run() {
	// 发送HTTP请求获取Node列表
	response := s.client.Get("/nodes/getAll")
	nodeKV := new(map[string]string)
	var nodeList []object.NodeStorage
	json.Unmarshal([]byte(response), nodeKV)
	for _, val := range *nodeKV {
		var node object.NodeStorage
		json.Unmarshal([]byte(val), &node)
		nodeList = append(nodeList, node)
	}
	s.nodeList = nodeList

	// 为目前的pod分配node
	s.scheduleAllPod()

	var wg sync.WaitGroup
	wg.Add(1)
	// 开始监听消息队列中pod的增量信息
	go func() {
		err := s.subscriber.SubscribeWithSync(s.podQueue, subscriber.Handler(s.podHandler), &wg)
		if err != nil {
			fmt.Printf(err.Error())
			_ = s.subscriber.CloseConnection()
			wg.Done()
		}
	}()

	// 开始监听node的增量信息
	go func() {
		err := s.subscriber.SubscribeWithSync(s.nodeQueue, subscriber.Handler(s.nodeHandler), &wg)
		if err != nil {
			fmt.Printf(err.Error())
			_ = s.subscriber.CloseConnection()
			wg.Done()
		}
	}()

	// 阻塞主协程，当至少有一个MQ停止监听时，主协程退出
	wg.Wait()
}

// --------------------POD STATUS LISTENER----------------
type podHandler struct {
	sched *Scheduler
}

func (h podHandler) Handle(jsonMsg []byte) {
	log.Println("scheduler get pod subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &podStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevPodStorage)
	switch msg.EventType {
	case object.CREATE:
		if podStorage.Node == "" {
			h.sched.roundRobin(&podStorage)
		}
	case object.UPDATE:
		if podStorage.Node == "" {
			h.sched.roundRobin(&podStorage)
		}
	case object.DELETE:
	}
}

// -------------------NODE STATUS LISTENER----------------
type nodeHandler struct {
	sched *Scheduler
}

func (h nodeHandler) Handle(jsonMsg []byte) {
	log.Println("Scheduler get node subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	nodeStorage := object.NodeStorage{}
	prevNodeStorage := object.NodeStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &nodeStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevNodeStorage)
	switch msg.EventType {
	case object.CREATE:
		h.sched.nodeList = append(h.sched.nodeList, nodeStorage)
		h.sched.scheduleAllPod()
	case object.UPDATE:
		var newList []object.NodeStorage
		for _, node := range h.sched.nodeList {
			if node.Node.Metadata.Uid == nodeStorage.Node.Metadata.Uid {
				newList = append(newList, nodeStorage)
			} else {
				newList = append(newList, node)
			}
		}
		h.sched.nodeList = newList
	case object.DELETE:
		var newList []object.NodeStorage
		for _, node := range h.sched.nodeList {
			if node.Node.Metadata.Uid == prevNodeStorage.Node.Metadata.Uid {
				continue
			} else {
				newList = append(newList, node)
			}
		}
		h.sched.nodeList = newList
		// 将分布在该node上的pod重新调配至其他node上
		h.sched.scheduleAllPod()
	}
}

// ----------------------SCHEDULE WORKER----------------------

var roundIndex = 0

func (s *Scheduler) roundRobin(pod *object.PodStorage) {
	nodeNum := len(s.nodeList)
	if nodeNum == 0 {
		if pod.Node == "" {
			return
		}
		pod.Node = ""
		updateMsg, _ := json.Marshal(pod)
		response := s.client.Post("/pods/update", updateMsg)
		if response != "ok" {
			fmt.Println("err when alloc node for pod " + pod.Config.Metadata.Name + " response: " + response)
		}
	} else {
		pod.Node = s.nodeList[roundIndex%nodeNum].Node.Metadata.Uid
		updateMsg, _ := json.Marshal(pod)
		response := s.client.Post("/pods/update", updateMsg)
		if response != "ok" {
			fmt.Println("err when alloc node for pod " + pod.Config.Metadata.Name + " response: " + response)
		}
	}
}

func (s *Scheduler) affinity() {

}

// 检查apiserver处的所有pod
func (s *Scheduler) scheduleAllPod() {
	// 发送HTTP请求获取Pod列表
	response := s.client.Get("/pods/getAll")
	podList := new(map[string]string)
	json.Unmarshal([]byte(response), podList)

	// 遍历pod列表，给没有节点的pod选择node节点
	log.Println("Len of PodList: " + strconv.Itoa(len(*podList)))
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		if !s.isNodeValid(podInfo.Node) {
			s.roundRobin(&podInfo)
		}
	}
}

//// 将原本跑在某一node上的pod全部重新分配（用于该node被停止或被删除）
//func (s *Scheduler) reschedule(nodeUid string) {
//	// 发送HTTP请求获取Pod列表
//	response := s.client.Get("/pods/getAll")
//	podList := new(map[string]string)
//	json.Unmarshal([]byte(response), podList)
//
//	// 遍历pod列表，给没有节点的pod选择node节点
//	log.Println("Len of PodList: " + string(len(*podList)))
//	for _, val := range *podList {
//		podInfo := object.PodStorage{}
//		_ = json.Unmarshal([]byte(val), &podInfo)
//		if podInfo.Node == "" || podInfo.Node == nodeUid {
//			s.roundRobin(&podInfo)
//		}
//	}
//}

func (s *Scheduler) isNodeValid(node string) bool {
	if node == "" {
		return false
	}
	for _, n := range s.nodeList {
		if n.Node.Metadata.Uid == node {
			return true
		}
	}
	return false
}
