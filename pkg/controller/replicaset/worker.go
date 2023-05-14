package replicaset

import (
	"k8s/object"
)

type Worker interface {
	Start()
	Stop()
	UpdateReplicaset(rs object.ReplicaSet)
	PodChangeHandler()
	SyncPods(rsPodNum int)
}

type worker struct {
	target object.ReplicaSet
	quit   chan int
}

func (w *worker) Start() {
	for {
		// watch(topic_pod, PodChangeHandler)

		select {
		case <-w.quit:
			return
		}
	}
}

func (w *worker) Stop() {
	w.quit <- 1
}

func (w *worker) UpdateReplicaset(rs object.ReplicaSet) {
	w.target = rs

	w.PodChangeHandler()
}

func (w *worker) PodChangeHandler() {
	// TODO 上锁

	// 设msg.pod为发生变化的Pod
	var msg_pod object.Pod

	if msg_pod.Metadata.Labels.App != w.target.Spec.Selector.MatchLabels.App ||
		msg_pod.Metadata.Labels.Env != w.target.Spec.Selector.MatchLabels.Env {
		return
	} else {
		// list(pods)
		var podsList []object.Pod
		rsPodNum := 0

		for _, value := range podsList {
			if value.Metadata.Labels.App == w.target.Spec.Selector.MatchLabels.App &&
				value.Metadata.Labels.Env == w.target.Spec.Selector.MatchLabels.Env {
				rsPodNum++
			}
		}

		w.SyncPods(rsPodNum)
	}
}

func (w *worker) SyncPods(rsPodNum int) {
	if rsPodNum < w.target.Spec.Replicas {
		// addPodsToApiserver
	} else if rsPodNum > w.target.Spec.Replicas {
		// deletePodsToApiserver
	}
}

func NewWorker(rs object.ReplicaSet, quit0 chan int) Worker {
	return &worker{
		target: rs,
		quit:   quit0,
	}
}
