package hpa

import (
	"k8s/object"
	"k8s/pkg/controllers/replicaset"
	"time"
)

type Worker interface {
	Start()
	Stop()
	UpdateHpa(hpa object.Hpa)

	//PodSyncHandler()

	//GetSelectedPodNum() (int, int, []int)
	//SyncPods()
}

type worker struct {
	target   object.Hpa
	RSworker replicaset.Worker
	cache    Cache

	//client通过http进行replicaset的增改删
	//client *HTTPClient.Client
}

func (w *worker) Start() {
	//每隔15s检查一次
	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {

		//log.Printf("[hpa worker] rs worker replica num: %d\n", w.RSworker.GetRS().Spec.Replicas)

		//podList := w.cache.GetPodStatusList()
		//
		//var cpuUtil, memUtil float64
		//cpuUtil = 0
		//memUtil = 0
		//rsPodNum := 0
		//for _, value := range podList {
		//	//fmt.Println(value)
		//	var pod object.PodStorage
		//	err := json.Unmarshal([]byte(value), &pod)
		//	if err != nil {
		//		log.Println("[hpa worker] unmarshall pod failed")
		//		return
		//	}
		//	if pod.Config.Metadata.Labels.App == w.RSworker.GetRS().Spec.Selector.MatchLabels.App &&
		//		pod.Config.Metadata.Labels.Env == w.RSworker.GetRS().Spec.Selector.MatchLabels.Env {
		//		rsPodNum++
		//		cpuUtil += pod.RunningMetrics.CPUUtil
		//		memUtil += pod.RunningMetrics.MemUtil
		//	}
		//}
		//cpuUtil /= (float64)(rsPodNum)
		//memUtil /= (float64)(rsPodNum)
		//log.Println("[hpa worker] hpa: "+w.target.Metadata.Name+", cpu util: %f, memory util: %f", cpuUtil, memUtil)
		//
		//num1 := math.Ceil(rsPodNum * (cpuUtil / w.target.Spec.Metrics))
	}
}

func (w *worker) Stop() {

}

func (w *worker) UpdateHpa(hpa object.Hpa) {

}

func NewWorker(hpa object.Hpa, cache Cache, RSworker replicaset.Worker) Worker {
	return &worker{
		target:   hpa,
		cache:    cache,
		RSworker: RSworker,
	}
}
