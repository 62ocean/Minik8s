package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/kubelet/pod"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"strconv"
	"sync"
	"time"
)

//----------------KUBELET---------------------

type Kubelet struct {
	client        *HTTPClient.Client
	node          object.Node
	podSubscriber *subscriber.Subscriber
	podQueue      string
	podHandler    podHandler
	pods          []object.PodStorage
	mutex         sync.Mutex
}

// NewKubelet kubelet对象的构造函数
func NewKubelet(name string) (*Kubelet, error) {
	// 使用HTTP，构建node对象传递到APIServer处
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	id, _ := uuid.NewUUID()
	nodeInfo := object.Node{
		Metadata: object.Metadata{
			Name:      name,
			Namespace: "default",
			Uid:       id.String(),
		},
		IP: "127.0.0.1",
	}
	info, _ := json.Marshal(nodeInfo)
	response := client.Post("/nodes/create", info)
	fmt.Println("get response from APIServer" + response)

	// 建立消息监听队列
	sub, _ := subscriber.NewSubscriber(global.MQHost)

	// 创建kubelet监听队列
	kub := Kubelet{
		client:        client,
		node:          nodeInfo,
		podSubscriber: sub,
		podQueue:      "pods_node",
	}
	h := podHandler{
		nodeID: kub.node.Metadata.Uid,
		kub:    &kub,
	}
	kub.podHandler = h
	return &kub, nil
}

// Run kubelet运行的入口函数
func (kub *Kubelet) Run() {

	// 发送HTTP请求获取Pod列表
	response := kub.client.Get("/pods/getAll")
	podList := new(map[string]string)
	json.Unmarshal([]byte(response), podList)

	// 遍历pod列表，运行在本node上的pod予以启动
	log.Println("Len of PodList: " + strconv.Itoa(len(*podList)))
	kub.mutex.Lock()
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		if podInfo.Node == kub.node.Metadata.Uid {
			kub.createPod(podInfo)
		}
	}
	kub.mutex.Unlock()

	//开启协程监听本地container变化 (每隔一秒轮询一次)
	go func() {
		// 等待初始化工作结束再开始监听
		time.Sleep(time.Second * 10)
		for {
			kub.watchPods()
			time.Sleep(time.Second * 3)
		}
	}()

	// 开始监听消息队列中pod的增量信息
	err := kub.podSubscriber.Subscribe(kub.podQueue, subscriber.Handler(kub.podHandler))
	if err != nil {
		fmt.Printf(err.Error())
		_ = kub.podSubscriber.CloseConnection()
	}
}

// StopKubelet 释放kubelet资源
func StopKubelet(kl *Kubelet) {
	if kl == nil {
		return
	}
	kl.client.Del("/nodes/remove/" + kl.node.IP)
}

// --------------------POD STATUS LISTENER----------------
type podHandler struct {
	nodeID string
	kub    *Kubelet
}

func (h podHandler) Handle(jsonMsg []byte) {
	h.kub.mutex.Lock()
	log.Println("Node get subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &podStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevPodStorage)
	switch msg.EventType {
	case object.CREATE:
		if podStorage.Node == h.nodeID {
			h.kub.createPod(podStorage)
		}
	case object.UPDATE:
		if prevPodStorage.Node == h.nodeID {
			if podStorage.Node != h.nodeID {
				// pod被转移至其他node
				h.kub.deletePod(podStorage)
			}
			if podStorage.Node == h.nodeID {
				if podStorage.Status == prevPodStorage.Status {
					// 对本node已有节点进行修改（若非单纯的状态变化，直接删除了pod重创即可）
					h.kub.deletePod(prevPodStorage)
					h.kub.createPod(podStorage)
					h.kub.watchPods()
				}
			}
		} else {
			// pod被转移至本node
			if podStorage.Node == h.nodeID {
				h.kub.createPod(podStorage)
			}
		}
	case object.DELETE:
		if prevPodStorage.Node == h.nodeID {
			h.kub.deletePod(prevPodStorage)
		}
	}
	h.kub.mutex.Unlock()
}

// ----------------------POD WORKER----------------------

func (kub *Kubelet) createPod(podInfo object.PodStorage) {
	//启动pod与相关容器
	log.Println("Begin to crate pod" + podInfo.Config.Metadata.Name)
	err := pod.CreatePod(&podInfo.Config)
	if err != nil {
		log.Println("Create pod error:")
		log.Println(err.Error())
		return
	}

	// 运行相关容器
	pod.StartPod(&podInfo.Config)

	//通知apiServer保存status
	podInfo.Status = object.RUNNING
	updateMsg, _ := json.Marshal(podInfo)
	resp := kub.client.Post("/pods/update", updateMsg)
	if resp == "ok" {
		log.Println("update pod's status to RUNNING")
		kub.pods = append(kub.pods, podInfo)
	} else {
		log.Println("cannot update pod's status after create")
	}
}

func (kub *Kubelet) deletePod(podInfo object.PodStorage) {
	log.Println("begin to delete pod" + podInfo.Config.Metadata.Name)
	//删除pod与相关容器
	pod.RemovePod(&podInfo.Config)
	var newPods []object.PodStorage
	for _, v := range kub.pods {
		if v.Config.Metadata.Uid != podInfo.Config.Metadata.Uid {
			newPods = append(newPods, v)
		}
	}
	kub.pods = newPods
}

// 定期监视本地容器,在不影响pod的情况下在自动重启容器（ps 双层循环嵌套，一个node管多个pod，一个pod有多个container）
func (kub *Kubelet) watchPods() {
	kub.mutex.Lock()
	for _, myPod := range kub.pods {
		update, err := pod.SyncPod(&myPod.Config)
		if err != nil {
			fmt.Println(err.Error())
			kub.mutex.Unlock()
			return
		}
		if update {
			updateMsg, _ := json.Marshal(myPod)
			resp := kub.client.Post("/pods/update", updateMsg)
			if resp == "ok" {
				log.Println("update pod's containers")
			} else {
				log.Println("cannot update pod's containers")
			}
		}
	}
	kub.mutex.Unlock()
}
