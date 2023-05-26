package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/kubelet/cache"
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
	pods          []*cache.PodCache
	toBeDel       []*cache.PodCache
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
	fmt.Println("get response from APIServer when register node: " + response)
	// 略等待scheduler分配pod
	time.Sleep(time.Millisecond * 500)

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

	// 开启协程定期更新同步podList（每30s同步一次）
	go func() {
		for {
			kub.syncPodList()
			time.Sleep(time.Second * 60)
		}
	}()

	//开启协程监听本地container变化，并更新上传pod状态 (每隔5s 轮询一次)
	go func() {
		time.Sleep(time.Millisecond * 500)
		for {
			kub.watchPods()
			time.Sleep(time.Second * 5)
		}
	}()

	// 开始监听消息队列中pod的增量信息
	////等待初始化工作结束再开始监听
	//time.Sleep(time.Second * 1)
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
	msg := object.MQMessage{}
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &podStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevPodStorage)
	switch msg.EventType {
	case object.CREATE:
		log.Println("Node get msg of CREATE")
		if podStorage.Node == h.nodeID {
			podCache := h.kub.addToList(podStorage)
			go h.kub.createPod(podCache)
		}
	case object.UPDATE:
		log.Println("Node get msg of UPDATE")
		// update 目前只可能是scheduler和状态同步两种情况造成，不用考虑太多
		if prevPodStorage.Node == h.nodeID {
			if podStorage.Node != h.nodeID {
				// pod被转移至其他node
				log.Println("Pod is removed to other node")
				podCache := h.kub.delFromList(podStorage)
				go h.kub.deletePod(*podCache)
			}
			if podStorage.Node == h.nodeID {
				log.Println("Pod status changed")
				// 此种情况只可能是状态同步，暂时不用管
				//if podStorage.Status == prevPodStorage.Status  {
				//	// 对本node已有节点进行修改
				//	h.kub.updateInList(podStorage)
				//}
			}
		} else {
			// pod被转移至本node
			log.Println("Pod is moved to this node")
			if podStorage.Node == h.nodeID {
				podCache := h.kub.addToList(podStorage)
				go h.kub.createPod(podCache)
			}
		}
	case object.DELETE:
		log.Println("Node get msg of DELETE")
		if prevPodStorage.Node == h.nodeID {
			podCache := h.kub.delFromList(prevPodStorage)
			go h.kub.deletePod(*podCache)
		}
	}
	h.kub.mutex.Unlock()
}

// 定期监视本地容器,在不影响pod的情况下在自动重启容器（ps 双层循环嵌套，一个node管多个pod，一个pod有多个container）
func (kub *Kubelet) watchPods() {
	kub.mutex.Lock()
	for _, delPod := range kub.toBeDel {
		log.Println("watchPods: delete pod")
		go kub.deletePod(*delPod)
	}
	kub.toBeDel = []*cache.PodCache{}
	for _, myPod := range kub.pods {
		if myPod.PodStorage.Status == object.STOPPED {
			// 未运行pod直接启动
			log.Println("watchPods: create pod")
			go kub.createPod(myPod)
		} else {
			// 已运行pod同步容器状态
			log.Println("watchPods: sync pod status")
			update, err := pod.SyncPod(myPod)
			if err != nil {
				fmt.Println(err.Error())
				kub.mutex.Unlock()
				return
			}
			if update {
				go func() {
					kub.deletePod(*myPod)
					kub.createPod(myPod)
				}()
			} else {
				go kub.uploadStatus(myPod.PodStorage)
			}
		}
	}
	kub.mutex.Unlock()
}

// 定期同步podList，上传pod资源消耗
func (kub *Kubelet) syncPodList() {
	kub.mutex.Lock()
	// 发送HTTP请求获取Pod列表
	response := kub.client.Get("/pods/getAll")
	podList := new(map[string]string)
	json.Unmarshal([]byte(response), podList)

	// 遍历pod列表，运行在本node上的加入podList
	log.Println("Len of PodList: " + strconv.Itoa(len(*podList)))
	var newList []*cache.PodCache
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		cacheInKub := kub.getPodCache(&podInfo)
		if cacheInKub == nil {
			podInfo.Status = object.STOPPED
			newCache := cache.PodCache{PodStorage: podInfo}
			newList = append(newList, &newCache)
		} else {
			newList = append(newList, cacheInKub)
		}
	}
	kub.pods = newList
	kub.mutex.Unlock()
}

// ----------------------POD WORKER----------------------

// 启动pod相关容器，填充cache中的容器id缓存
func (kub *Kubelet) createPod(podInfo *cache.PodCache) {
	//启动pod与相关容器
	log.Println("Begin to crate pod " + podInfo.PodStorage.Config.Metadata.Name)
	containers, err := pod.CreatePod(podInfo.PodStorage.Config)
	if err != nil {
		log.Println("Create pod error:")
		log.Println(err.Error())
		return
	}
	// 运行相关容器
	pod.StartPod(containers)

	//通知apiServer保存status和资源利用率
	matrix := pod.GetStatusOfPod(podInfo)
	podInfo.PodStorage.RunningMetrics = matrix
	podInfo.PodStorage.Status = object.RUNNING
	kub.uploadStatus(podInfo.PodStorage)
	podInfo.ContainerMeta = containers
}

func (kub *Kubelet) deletePod(podInfo cache.PodCache) {
	log.Println("begin to delete pod" + podInfo.PodStorage.Config.Metadata.Name)
	//删除pod与相关容器
	pod.RemovePod(&podInfo)
}

// -----------------------------TOOLS---------------------------

func (kub *Kubelet) getPodCache(storage *object.PodStorage) *cache.PodCache {
	for _, podCache := range kub.pods {
		if podCache.PodStorage.Config.Metadata.Uid == storage.Config.Metadata.Uid {
			return podCache
		}
	}
	return nil
}

func (kub *Kubelet) addToList(storage object.PodStorage) *cache.PodCache {
	// 不存在时添加
	//if kub.getPodCache(&storage) == nil {
	log.Println("add pod into podList")
	storage.Status = object.STOPPED
	newCache := cache.PodCache{PodStorage: storage}
	kub.pods = append(kub.pods, &newCache)
	return &newCache
	//}
}

func (kub *Kubelet) delFromList(storage object.PodStorage) *cache.PodCache {
	// 存在时删除
	//if kub.getPodCache(&storage) != nil {
	log.Println("del pod from podList")
	var newPods []*cache.PodCache
	var ret *cache.PodCache
	for _, v := range kub.pods {
		if v.PodStorage.Config.Metadata.Uid != storage.Config.Metadata.Uid {
			newPods = append(newPods, v)
		} else {
			//var copyOfPod cache.PodCache
			//// 深拷贝
			//temp, _ := json.Marshal(v)
			//_ = json.Unmarshal(temp, &copyOfPod)
			kub.toBeDel = append(kub.toBeDel, v)
			ret = v
		}
	}
	kub.pods = newPods
	return ret

	//}
}

func (kub *Kubelet) updateInList(storage object.PodStorage) {
	podInList := kub.getPodCache(&storage)
	if podInList != nil {
		// 存在时修改
		log.Println("update pod in podList")
		podInList.PodStorage = storage
		podInList.PodStorage.Status = object.STOPPED
	} else {
		// 不存在时添加
		kub.addToList(storage)
	}
}

func (kub *Kubelet) uploadStatus(podInfo object.PodStorage) {
	updateMsg, _ := json.Marshal(podInfo)
	resp := kub.client.Post("/pods/update", updateMsg)
	if resp == "ok" {
		log.Println("upload pod's status")
	} else {
		log.Println("cannot upload pod's status, response: " + resp)
	}
}
