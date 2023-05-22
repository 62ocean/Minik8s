package worker

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"strconv"
)

type Worker interface {
	Start()
	Stop()
	UpdateReplicaset(rs object.ReplicaSet)
	//PodSyncHandler()

	GetSelectedPodNum() (int, int)
	SyncPods()
}

type worker struct {
	target object.ReplicaSet
	quit   chan int

	//s监听pod的变化，handler处理
	s       *subscriber.Subscriber
	handler *PodSyncHandler

	//client通过http进行replicaset的增改删
	client *HTTPClient.Client
}

func (w *worker) Start() {
	//for {
	//	// watch(topic_pod, PodChangeHandler)
	//
	//	select {
	//	case <-w.quit:
	//		return
	//	}
	//}
	fmt.Println("worker start")

	//创建client对pod进行增删改操作
	w.client = HTTPClient.CreateHTTPClient(global.ServerHost)
	//worker启动时先检查一下pod数量是否符合要求
	w.SyncPods()

	//创建subscribe监听pod的变化
	w.s, _ = subscriber.NewSubscriber("amqp://guest:guest@localhost:5672/")
	w.handler = NewPodSyncHandler(w)
	err := w.s.Subscribe("pods", subscriber.Handler(w.handler))
	if err != nil {
		fmt.Println("subcribe pods failed")
		return
	}
}

func (w *worker) Stop() {
	w.quit <- 1
}

func (w *worker) UpdateReplicaset(rs object.ReplicaSet) {
	w.target = rs

	//w.PodChangeHandler()
}

func (w *worker) GetSelectedPodNum() (int, int) {
	//得到所有的pod列表
	response := w.client.Get("/pods/getAll")
	podList := new(map[string]string)
	err := json.Unmarshal([]byte(response), podList)
	if err != nil {
		fmt.Println("unmarshall podlist failed")
		return -1, -1
	}

	//fmt.Println(podList)

	// 统计符合要求的pod个数
	num := 0
	maxRepIndex := 0
	for _, value := range *podList {
		fmt.Println(value)
		var pod object.PodStorage
		err := json.Unmarshal([]byte(value), &pod)
		if err != nil {
			fmt.Println("unmarshall pod failed")
			return -1, -1
		}
		if pod.Config.Metadata.Labels.App == w.target.Spec.Selector.MatchLabels.App &&
			pod.Config.Metadata.Labels.Env == w.target.Spec.Selector.MatchLabels.Env {
			num++
			if pod.Replica > maxRepIndex {
				maxRepIndex = pod.Replica
			}
		}
	}

	fmt.Println(num)

	return num, maxRepIndex
}

//func (w *worker) PodSyncHandler(pod object.Pod) {
//	// TODO 上锁
//
//	// 设msg.pod为发生变化的Pod
//
//	if pod.Metadata.Labels.App != w.target.Spec.Selector.MatchLabels.App ||
//		pod.Metadata.Labels.Env != w.target.Spec.Selector.MatchLabels.Env {
//		return
//	} else {
//		// list(pods)
//		//var podsList []object.Pod
//		//rsPodNum := 0
//		//
//		//for _, value := range podsList {
//		//	if value.Metadata.Labels.App == w.target.Spec.Selector.MatchLabels.App &&
//		//		value.Metadata.Labels.Env == w.target.Spec.Selector.MatchLabels.Env {
//		//		rsPodNum++
//		//	}
//		//}
//
//		w.SyncPods(w.GetSelectedPodNum())
//	}
//}

func (w *worker) SyncPods() {
	podTemplate := w.target.Spec.PodTemplate

	rsPodNum, maxRepIndex := w.GetSelectedPodNum()
	for rsPodNum != w.target.Spec.Replicas {
		if rsPodNum < w.target.Spec.Replicas {
			// 修改pod uid，名字以及容器名称 (ps要用深拷贝，防止修改podTemplate)
			temp := &podTemplate
			newPod := *temp
			maxRepIndex = maxRepIndex + 1
			id, _ := uuid.NewUUID()
			newPod.Metadata.Uid = id.String()
			newPod.Metadata.Name = podTemplate.Metadata.Name + "-" + strconv.Itoa(maxRepIndex)
			oldContainers := podTemplate.Spec.Containers
			var newContainers []object.Container
			for _, c := range oldContainers {
				newC := c
				newC.Name = c.Name + "-" + strconv.Itoa(maxRepIndex)
				newContainers = append(newContainers, newC)
			}
			newPod.Spec.Containers = newContainers
			var podJson []byte
			podJson, _ = json.Marshal(newPod)

			w.client.Post("/pods/create", podJson)

			rsPodNum++

		} else if rsPodNum > w.target.Spec.Replicas {
			// deletePodsToApiserver
		}
	}

}

func NewWorker(rs object.ReplicaSet, quit0 chan int) Worker {
	return &worker{
		target: rs,
		quit:   quit0,
	}
}
