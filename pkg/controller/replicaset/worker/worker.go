package worker

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"strconv"
)

type Worker interface {
	Start()
	Stop()
	UpdateReplicaset(rs object.ReplicaSet)
	//PodSyncHandler()

	GetSelectedPodNum() (int, int, []int)
	SyncPods()
}

type worker struct {
	target object.ReplicaSet

	//s监听pod的变化，handler处理
	s       *subscriber.Subscriber
	handler *PodSyncHandler

	//client通过http进行replicaset的增改删
	client *HTTPClient.Client
}

func (w *worker) Start() {

	log.Println("worker start")

	//创建client对pod进行增删改操作
	w.client = HTTPClient.CreateHTTPClient(global.ServerHost)
	//worker启动时先检查一下pod数量是否符合要求
	w.SyncPods()

	//创建subscribe监听pod的变化
	w.s, _ = subscriber.NewSubscriber("amqp://guest:guest@localhost:5672/")
	w.handler = NewPodSyncHandler(w)
	err := w.s.Subscribe("pods_"+w.target.Spec.Selector.MatchLabels.App, subscriber.Handler(w.handler))
	if err != nil {
		fmt.Println("subcribe pods failed")
		return
	}
}

func (w *worker) Stop() {
	err := w.s.CloseConnection()
	if err != nil {
		fmt.Println("close connection error")
		return
	}
}

func (w *worker) UpdateReplicaset(rs object.ReplicaSet) {
	w.target = rs

	w.SyncPods()
}

func (w *worker) GetSelectedPodNum() (int, int, []int) {
	//得到所有的pod列表
	response := w.client.Get("/pods/getAll")
	podList := new(map[string]string)
	err := json.Unmarshal([]byte(response), podList)
	if err != nil {
		fmt.Println("unmarshall podlist failed")
		return -1, -1, nil
	}

	//log.Println(podList)

	// 统计符合要求的pod个数
	num := 0
	maxRepIndex := 0
	var seqNum []int
	for _, value := range *podList {
		//fmt.Println(value)
		var pod object.PodStorage
		err := json.Unmarshal([]byte(value), &pod)
		if err != nil {
			fmt.Println("unmarshall pod failed")
			return -1, -1, nil
		}
		if pod.Config.Metadata.Labels.App == w.target.Spec.Selector.MatchLabels.App &&
			pod.Config.Metadata.Labels.Env == w.target.Spec.Selector.MatchLabels.Env {
			num++
			seqNum = append(seqNum, pod.Replica)
			if pod.Replica > maxRepIndex {
				maxRepIndex = pod.Replica
			}
		}
	}

	log.Println(num)

	return num, maxRepIndex, seqNum
}

func (w *worker) SyncPods() {
	podTemplate := w.target.Spec.PodTemplate

	rsPodNum, maxRepIndex, seqNum := w.GetSelectedPodNum()
	for rsPodNum != w.target.Spec.Replicas {
		if rsPodNum < w.target.Spec.Replicas {
			// 修改pod uid，名字以及容器名称 (ps要用深拷贝，防止修改podTemplate)
			temp := &podTemplate
			newPod := *temp
			maxRepIndex = maxRepIndex + 1
			id, _ := uuid.NewUUID()
			newPod.Metadata.Uid = id.String()
			newPod.Metadata.Name = podTemplate.Metadata.Name + "-" + strconv.Itoa(maxRepIndex)
			var podJson []byte
			podJson, _ = json.Marshal(newPod)

			w.client.Post("/pods/create", podJson)

			rsPodNum++

		} else if rsPodNum > w.target.Spec.Replicas {
			rmPodName := podTemplate.Metadata.Name + "-" + strconv.Itoa(seqNum[0])
			fmt.Println("remove seq num: " + strconv.Itoa(seqNum[0]))
			seqNum = seqNum[1:]

			var podJson []byte
			podJson, _ = json.Marshal(rmPodName)

			w.client.Post("/pods/remove", podJson)

			rsPodNum--
		}
	}

}

func NewWorker(rs object.ReplicaSet) Worker {
	return &worker{
		target: rs,
	}
}
