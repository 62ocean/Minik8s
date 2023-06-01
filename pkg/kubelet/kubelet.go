package kubelet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/kubelet/cache"
	"k8s/pkg/kubelet/pod"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"os"
	"strconv"
	"strings"
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
	pods          map[string]*cache.PodCache
	toBeDel       map[string]*cache.PodCache
	mutex         sync.Mutex
}

// NewKubelet kubelet对象的构造函数
func NewKubelet(name string) (*Kubelet, error) {
	// 使用HTTP，构建node对象传递到APIServer处
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	config := readConfiguration("./build/flannel.properties")
	id, _ := uuid.NewUUID()
	nodeInfo := object.Node{
		Metadata: object.Metadata{
			Name:      name,
			Namespace: "default",
			Uid:       id.String(),
		},
		IP: config["node-ip"],
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
		pods:          make(map[string]*cache.PodCache),
		toBeDel:       make(map[string]*cache.PodCache),
	}
	h := podHandler{
		nodeIP: kub.node.IP,
		kub:    &kub,
	}
	kub.podHandler = h
	return &kub, nil
}

// Run kubelet运行的入口函数
func (kub *Kubelet) Run() {
	// 优化：初始开启容器还是得快，直接拉列表并创建
	// 发送HTTP请求获取Pod列表
	response := kub.client.Get("/pods/getAll")
	podList := new(map[string]string)
	json.Unmarshal([]byte(response), podList)

	// 遍历pod列表，运行在本node上的pod予以启动
	log.Println("Len of PodList: " + strconv.Itoa(len(*podList)))
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		if podInfo.Node == kub.node.IP {
			podCache := kub.addToList(podInfo)
			podCache.PodStorage.Status = object.PENDING
			go kub.createPod(podCache)
		}
	}

	// 开启协程定期更新同步podList（每1min同步一次）
	go func() {
		for {
			time.Sleep(time.Second * 60)
			kub.syncPodList()
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
	nodeIP string
	kub    *Kubelet
}

func (h podHandler) Handle(jsonMsg []byte) {
	msg := object.MQMessage{}
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &podStorage)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevPodStorage)
	h.kub.mutex.Lock()
	switch msg.EventType {
	case object.CREATE:
		log.Println("Node get msg of CREATE")
		if podStorage.Node == h.nodeIP {
			podCache := h.kub.addToList(podStorage)
			podCache.PodStorage.Status = object.PENDING
			go h.kub.createPod(podCache)
		}
	case object.UPDATE:
		log.Println("Node get msg of UPDATE")
		// update 目前只可能是scheduler和状态同步两种情况造成，不用考虑太多
		if prevPodStorage.Node == h.nodeIP {
			if podStorage.Node != h.nodeIP {
				// pod被转移至其他node
				log.Println("Pod is removed to other node")
				podCache := h.kub.DelFromList(podStorage)
				go h.kub.deletePod(*podCache)
			}
			if podStorage.Node == h.nodeIP {
				log.Println("Pod status changed")
				// 此种情况只可能是状态同步、或者分配pod IP，暂时不用管
				//if podStorage.Status == prevPodStorage.Status  {
				//	// 对本node已有节点进行修改
				//	h.kub.updateInList(podStorage)
				//}
			}
		} else {
			// pod被转移至本node
			log.Println("Pod is moved to this node")
			if podStorage.Node == h.nodeIP {
				podCache := h.kub.addToList(podStorage)
				podCache.PodStorage.Status = object.PENDING
				go h.kub.createPod(podCache)
			}
		}
	case object.DELETE:
		log.Println("Node get msg of DELETE")
		if prevPodStorage.Node == h.nodeIP {
			podCache := h.kub.DelFromList(prevPodStorage)
			if podCache == nil {
				podCache = h.kub.toBeDel[prevPodStorage.Config.Metadata.Uid]
			}
			go h.kub.deletePod(*podCache)
		}
	}
	h.kub.mutex.Unlock()
}

// 定期监视本地容器,在不影响pod的情况下在自动重启容器（ps 双层循环嵌套，一个node管多个pod，一个pod有多个container）
func (kub *Kubelet) watchPods() {
	kub.mutex.Lock()
	for _, delPod := range kub.toBeDel {
		log.Println("WatchPods: delete pod " + delPod.PodStorage.Config.Metadata.Name)
		go kub.deletePod(*delPod)
	}
	kub.toBeDel = make(map[string]*cache.PodCache)
	for _, myPod := range kub.pods {
		if myPod.PodStorage.Status == object.STOPPED {
			// 未运行pod直接启动
			log.Println("WatchPods: create pod " + myPod.PodStorage.Config.Metadata.Name)
			go kub.createPod(myPod)
		} else if myPod.PodStorage.Status == object.RUNNING {
			// 已运行pod同步容器状态
			log.Println("WatchPods: sync pod " + myPod.PodStorage.Config.Metadata.Name)
			update := pod.SyncPod(myPod)
			if update {
				myPod.PodStorage.Status = object.PENDING
				tempPod := myPod
				go func() {
					kub.deletePod(*tempPod)
					kub.createPod(tempPod)
				}()
			} else {
				tempPod := myPod
				go kub.uploadStatus(tempPod.PodStorage)
			}
		}
		// pending 状态说明正在重启，暂时停止上传状态和核对状态
	}
	kub.mutex.Unlock()
}

// 定期同步podList，上传pod资源消耗
func (kub *Kubelet) syncPodList() {
	kub.mutex.Lock()
	log.Println("Begin to sync pod List")
	// 发送HTTP请求获取Pod列表
	response := kub.client.Get("/pods/getAll")
	podList := new(map[string]string)
	podMap := make(map[string]*object.PodStorage)
	json.Unmarshal([]byte(response), podList)

	// 遍历pod列表，运行在本node上的加入podList
	log.Println("Len of PodList: " + strconv.Itoa(len(*podList)))
	podListLen := 0
	for _, val := range *podList {
		podInfo := object.PodStorage{}
		_ = json.Unmarshal([]byte(val), &podInfo)
		if podInfo.Node == kub.node.IP {
			podListLen++
			kub.addToList(podInfo)
			podMap[podInfo.Config.Metadata.Uid] = &podInfo
		}
	}
	// 核对长度，不对劲说明有要删除的内容 (ps.map遍历时删除是安全的可以放心)
	if len(kub.pods) > podListLen {
		for key, _ := range kub.pods {
			if podMap[key] == nil {
				kub.moveToDelList(kub.pods[key].PodStorage)
			}
		}
	}
	kub.mutex.Unlock()
}

// ----------------------POD WORKER----------------------

// 启动pod相关容器，填充cache中的容器id缓存
func (kub *Kubelet) createPod(podInfo *cache.PodCache) {
	//启动pod与相关容器
	log.Println("Begin to crate pod " + podInfo.PodStorage.Config.Metadata.Name)
	containers, err := pod.CreatePod(&podInfo.PodStorage.Config)
	podInfo.ContainerMeta = containers
	if err != nil {
		log.Println("Create pod error:")
		log.Println(err.Error())
		return
	}
	// 运行相关容器
	pod.StartPod(containers, podInfo.PodStorage.Config.Metadata.Name)

	//通知apiServer保存status和资源利用率
	matrix := pod.GetStatusOfPod(podInfo)
	podInfo.PodStorage.RunningMetrics = matrix
	podInfo.PodStorage.Status = object.RUNNING
	kub.uploadStatus(podInfo.PodStorage)
}

func (kub *Kubelet) deletePod(podInfo cache.PodCache) {
	log.Println("begin to delete pod " + podInfo.PodStorage.Config.Metadata.Name)
	//删除pod与相关容器
	pod.RemovePod(&podInfo)
}

// -----------------------------TOOLS---------------------------

func (kub *Kubelet) getPodCache(storage *object.PodStorage) *cache.PodCache {
	val, ok := kub.pods[storage.Config.Metadata.Uid]
	if ok {
		return val
	} else {
		return nil
	}
}

func (kub *Kubelet) addToList(storage object.PodStorage) *cache.PodCache {
	val := kub.getPodCache(&storage)
	if val == nil {
		log.Println("add pod into podList")
		storage.Status = object.STOPPED
		newCache := cache.PodCache{PodStorage: storage}
		kub.pods[storage.Config.Metadata.Uid] = &newCache
		return &newCache
	} else {
		return val
	}
}

func (kub *Kubelet) DelFromList(storage object.PodStorage) *cache.PodCache {
	log.Println("del pod from podList")
	key := storage.Config.Metadata.Uid
	oldVal := kub.pods[key]
	delete(kub.pods, key)
	return oldVal
}

func (kub *Kubelet) moveToDelList(storage object.PodStorage) *cache.PodCache {
	log.Println("move pod ", storage.Config.Metadata.Name, " from podList to delList")
	key := storage.Config.Metadata.Uid
	oldVal := kub.pods[key]
	delete(kub.pods, key)
	kub.toBeDel[key] = oldVal
	return oldVal
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
	//log.Print("update status of pod: ")
	//log.Println(podInfo)
	updateMsg, _ := json.Marshal(podInfo)
	resp := kub.client.Post("/pods/update", updateMsg)
	if resp == "ok" {
		log.Println("upload pod's status")
	} else {
		log.Println("cannot upload pod's status, response: " + resp)
	}
}

func readConfiguration(configurationFile string) map[string]string {
	var properties = make(map[string]string)
	confFile, err := os.OpenFile(configurationFile, os.O_RDONLY, 0666)
	defer func(confFile *os.File) {
		if err := confFile.Close(); err != nil {
			panic(err)
		}
	}(confFile)
	if err != nil {
		fmt.Printf("The config file %s is not exits.", configurationFile)
	} else {
		reader := bufio.NewReader(confFile)
		for {
			if confString, err := reader.ReadString('\n'); err != nil {
				if err == io.EOF {
					break
				}
			} else {
				if len(confString) == 0 || confString == "\n" || confString[0] == '#' {
					continue
				}
				properties[strings.Split(confString, "=")[0]] = strings.Replace(strings.Split(confString, "=")[1], "\n", "", -1)
			}
		}
	}
	return properties
}
