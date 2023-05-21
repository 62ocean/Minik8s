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
)

//----------------KUBELET---------------------

type Kubelet struct {
	client        *HTTPClient.Client
	node          object.Node
	podSubscriber *subscriber.Subscriber
	podQueue      string
	podHandler    podHandler
	pods          []object.PodStorage
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
	log.Println("Len of PodList: " + string(len(*podList)))
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		if podInfo.Node == kub.node.Metadata.Uid {
			kub.createPod(podInfo)
		}
	}

	// 开始监听消息队列中pod的增量信息
	err := kub.podSubscriber.Subscribe(kub.podQueue, subscriber.Handler(kub.podHandler))
	if err != nil {
		fmt.Printf(err.Error())
		_ = kub.podSubscriber.CloseConnection()
	}
}

// --------------------POD STATUS LISTENER----------------
type podHandler struct {
	nodeID string
	kub    *Kubelet
}

func (h podHandler) Handle(jsonMsg []byte) {
	log.Println("Node get subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &podStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevPodStorage)
	log.Println("type： " + string(rune(msg.EventType)))
	switch msg.EventType {
	case object.CREATE:
		if podStorage.Node == h.nodeID {
			h.kub.createPod(podStorage)
		}
		break
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
				}
			}
		} else {
			// pod被转移至本node
			if podStorage.Node == h.nodeID {
				h.kub.createPod(podStorage)
			}
		}
		break
	case object.DELETE:
		if prevPodStorage.Node == h.nodeID {
			h.kub.deletePod(prevPodStorage)
		}
		break
	}
}

// ----------------------POD WORKER----------------------
func (kub Kubelet) createPod(podInfo object.PodStorage) {
	//启动pod与相关容器
	log.Println("begin to crate pod" + podInfo.Config.Metadata.Name)
	err := pod.CreatePod(&podInfo.Config)
	if err != nil {
		log.Println("Create pod error:")
		log.Println(err.Error())
		return
	}

	// 运行相关容器
	err = pod.StartPod(&podInfo.Config)
	if err != nil {
		log.Println("Create pod error:")
		log.Println(err.Error())
		return
	}

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

func (kub Kubelet) deletePod(podInfo object.PodStorage) {
	log.Println("begin to delete pod" + podInfo.Config.Metadata.Name)
	//删除pod与相关容器
	err := pod.RemovePod(&podInfo.Config)
	if err != nil {
		log.Println("Remove pod error:")
		log.Println(err.Error())
		return
	}
	//通知apiServer保存status(不保存，防止数据竞争)————scheduler要在改位置的时候把STATUS设置为STOPPED
	//podInfo.Status = object.STOPPED
	//updateMsg, _ := json.Marshal(podInfo)
	//resp := kub.client.Post("/pods/update", updateMsg)
	//if resp == "ok" {
	var newPods []object.PodStorage
	for _, v := range kub.pods {
		if v.Config.Metadata.Uid != podInfo.Config.Metadata.Uid {
			newPods = append(newPods, v)
		}
	}
	kub.pods = newPods

}
