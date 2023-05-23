package hpa

import (
	"k8s/object"
	"k8s/pkg/util/HTTPClient"
)

type Worker interface {
	Start()
	Stop()
	GetCPUMetric()
	//UpdateReplicaset(rs object.Hpa)
	//PodSyncHandler()

	//GetSelectedPodNum() (int, int, []int)
	//SyncPods()
}

type worker struct {
	target object.Hpa

	//client通过http进行replicaset的增改删
	client *HTTPClient.Client
}
