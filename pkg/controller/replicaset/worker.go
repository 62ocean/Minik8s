package replicaset

import "k8s/pkg/api/pod"

type Worker interface {
	Start()
	Stop()
	UpdateReplicaset(rs ReplicaSet)
	PodChangeHandler()
	SyncPods(rsPodNum int)
}

type worker struct {
	target ReplicaSet
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

func (w *worker) UpdateReplicaset(rs ReplicaSet) {
	w.target = rs
}

func (w *worker) PodChangeHandler() {
	// 设msg.pod为发生变化的Pod
	var msg_pod pod.Pod

	if msg_pod.Metadata.Labels.App != w.target.Spec.Selector.MatchLabels.App ||
		msg_pod.Metadata.Labels.Env != w.target.Spec.Selector.MatchLabels.Env {
		return
	} else {
		// list(pods)
		var podsList []pod.Pod
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

func NewWorker(rs ReplicaSet, quit0 chan int) Worker {
	return &worker{
		target: rs,
		quit:   quit0,
	}
}
