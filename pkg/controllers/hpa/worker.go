package hpa

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/controllers/replicaset"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"log"
	"math"
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
	ticker   *time.Ticker

	//client通过http进行replicaset的增改删
	client *HTTPClient.Client
}

func (w *worker) Start() {
	//每隔15s检查一次
	w.ticker = time.NewTicker(time.Second * 15)
	w.client = HTTPClient.CreateHTTPClient(global.ServerHost)

	for range w.ticker.C {

		log.Printf("[hpa worker] hpa: "+w.target.Metadata.Name+", rs pod num: %d\n", w.RSworker.GetRS().Spec.Replicas)

		podList := w.cache.GetPodStatusList()

		var cpuUtil, memUtil float64
		cpuUtil = 0
		memUtil = 0
		rsPodNum := 0
		for _, value := range podList {
			//fmt.Println(value)
			var pod object.PodStorage
			err := json.Unmarshal([]byte(value), &pod)
			if err != nil {
				log.Println("[hpa worker] unmarshall pod failed")
				return
			}
			if pod.Config.Metadata.Labels.App == w.RSworker.GetRS().Spec.Selector.MatchLabels.App &&
				pod.Config.Metadata.Labels.Env == w.RSworker.GetRS().Spec.Selector.MatchLabels.Env {
				rsPodNum++
				cpuUtil += pod.RunningMetrics.CPUUtil
				memUtil += pod.RunningMetrics.MemUtil
			}
		}
		cpuUtil /= (float64)(rsPodNum)
		memUtil /= (float64)(rsPodNum)
		log.Printf("[hpa worker] hpa: %s, cpu average util: %f, memory average util: %f\n", w.target.Metadata.Name, cpuUtil, memUtil)

		var cpuMetric, memMetric float64
		// 只考虑了cpu和memory两种指标
		for _, value := range w.target.Spec.Metrics {
			if value.Resource.Name == "cpu" {
				cpuMetric = value.Resource.Target.AverageUtilization
			} else {
				memMetric = value.Resource.Target.AverageUtilization
			}
		}

		// 计算两种指标的期望副本数
		num1 := (int)(math.Ceil((float64)(rsPodNum) * (cpuUtil / cpuMetric)))
		num2 := (int)(math.Ceil((float64)(rsPodNum) * (memUtil / memMetric)))
		log.Printf("[hpa worker] hpa: %s,(cpu) expect pod num: %d\n", w.target.Metadata.Name, num1)
		log.Printf("[hpa worker] hpa: %s,(memory) expect pod num: %d\n", w.target.Metadata.Name, num2)

		//HorizontalPodAutoscaler 采用为每个指标推荐的最大比例，并将工作负载设置为该大小（前提是这不大于你配置的总体最大值）。
		if (num1 > rsPodNum || num2 > rsPodNum) && rsPodNum < w.target.Spec.MaxReplicas {
			// 增加replica数量 (+1)
			newrs := w.RSworker.GetRS()
			newrs.Spec.Replicas++
			var rsJson []byte
			rsJson, _ = json.Marshal(newrs)

			w.client.Post("/replicasets/update", rsJson)

			log.Printf("[hpa worker] hpa: %s, add replica\n", w.target.Metadata.Name)
		} else if num1 < rsPodNum && num2 < rsPodNum && rsPodNum > w.target.Spec.MinReplicas {
			// 减少replica数量 (-1)
			newrs := w.RSworker.GetRS()
			newrs.Spec.Replicas--
			var rsJson []byte
			rsJson, _ = json.Marshal(newrs)

			w.client.Post("/replicasets/update", rsJson)

			log.Printf("[hpa worker] hpa: %s, minus replica\n", w.target.Metadata.Name)
		}
	}
}

func (w *worker) Stop() {
	w.ticker.Stop()
}

func (w *worker) UpdateHpa(hpa object.Hpa) {
	w.target = hpa

	// 暂时不考虑hpa的目标rs发生改变的情况
}

func NewWorker(hpa object.Hpa, cache Cache, RSworker replicaset.Worker) Worker {
	return &worker{
		target:   hpa,
		cache:    cache,
		RSworker: RSworker,
	}
}
