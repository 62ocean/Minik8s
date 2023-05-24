package hpa

import (
	"k8s/object"
)

type Worker interface {
	Start()
	Stop()
	UpdateHpa(hpa object.Hpa)
	GetCPUMetric()

	//PodSyncHandler()

	//GetSelectedPodNum() (int, int, []int)
	//SyncPods()
}

type worker struct {
	target object.Hpa

	//client通过http进行replicaset的增改删
	//client *HTTPClient.Client
}

func (w *worker) Start() {

}

func (w *worker) Stop() {

}

func (w *worker) UpdateHpa(hpa object.Hpa) {

}

func (w *worker) GetCPUMetric() {

}

func NewWorker(hpa object.Hpa) Worker {
	return &worker{
		target: hpa,
	}
}
