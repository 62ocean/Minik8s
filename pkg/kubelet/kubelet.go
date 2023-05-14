package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/kubelet/pod"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"os"
)

const serverHost = "http://127.0.0.1:8080"

type Kubelet struct {
	client        *HTTPClient.Client
	node          object.Node
	podSubscriber *subscriber.Subscriber
	podQueue      string
	podHandler    podHandler
	pods          []object.PodStorage
}

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
				// 对本node已有节点进行修改（直接删除了pod重创即可）
				h.kub.deletePod(prevPodStorage)
				h.kub.createPod(podStorage)
			}
		} else {
			// pod被转移至本node
			if podStorage.Node == h.nodeID {
				h.kub.createPod(podStorage)
			}
		}
		break
	case object.DELETE:
		if podStorage.Node == h.nodeID {
			h.kub.deletePod(podStorage)
		}
		break
	}
}

// Run kubelet运行的入口函数
func (kub *Kubelet) Run() {
	// -------------------FOR TEST------------------------
	// 解析pod的yaml配置文件
	dataBytes, err := os.ReadFile("pkg/kubelet/pod/podConfigTest.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	var podData object.Pod
	err2 := yaml.Unmarshal(dataBytes, &podData)
	if err2 != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	id, _ := uuid.NewUUID()
	fmt.Println("new pod uid" + id.String())
	podData.Metadata.Uid = id.String()
	fmt.Println(podData)
	podJson, _ := json.Marshal(podData)
	kub.client.Post("/pods/create", podJson)
	// -------------------FOR TEST------------------------

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
	err = kub.podSubscriber.Subscribe(kub.podQueue, subscriber.Handler(kub.podHandler))
	if err != nil {
		fmt.Printf(err.Error())
		_ = kub.podSubscriber.CloseConnection()
	}
}

// NewKubelet kubelet对象的构造函数
func NewKubelet(name string) (*Kubelet, error) {
	// 使用HTTP，构建node对象传递到APIServer处
	client := HTTPClient.CreateHTTPClient(serverHost)
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
		nodeID: id.String(),
		kub:    &kub,
	}
	kub.podHandler = h
	return &kub, nil
}

func (kub Kubelet) createPod(podInfo object.PodStorage) {
	//启动pod与相关容器
	log.Println("begin to crate pod" + podInfo.Config.Metadata.Name)
	err := pod.CreatePod(podInfo.Config)
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
		kub.pods = append(kub.pods, podInfo)
	}
}

func (kub Kubelet) deletePod(podInfo object.PodStorage) {
	log.Println("begin to delete pod" + podInfo.Config.Metadata.Name)
	//删除pod与相关容器
	err := pod.RemovePod(podInfo.Config)
	if err != nil {
		log.Println("Remove pod error:")
		log.Println(err.Error())
		return
	}

	//通知apiServer保存status
	podInfo.Status = object.STOPPED
	updateMsg, _ := json.Marshal(podInfo)
	resp := kub.client.Post("/pods/update", updateMsg)
	if resp == "ok" {
		var newPods []object.PodStorage
		for _, v := range kub.pods {
			if v.Config.Metadata.Uid != podInfo.Config.Metadata.Uid {
				newPods = append(newPods, v)
			}
		}
		kub.pods = newPods
	}
}
