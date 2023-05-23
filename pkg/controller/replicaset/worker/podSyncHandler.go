package worker

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"log"
)

type PodSyncHandler struct {
	w *worker
}

func (h *PodSyncHandler) Handle(msg []byte) {

	log.Println("pod receive msg: " + string(msg))

	var msgObject object.MQMessage
	err := json.Unmarshal(msg, &msgObject)
	if err != nil {
		fmt.Println("[worker] unmarshall msg failed")
		return
	}

	//var pod object.Pod
	//err = json.Unmarshal([]byte(msgObject.Value), &pod)
	//if err != nil {
	//	fmt.Println("[worker] unmarshall changed pod failed")
	//	return
	//}

	h.w.SyncPods()
}

func NewPodSyncHandler(w *worker) *PodSyncHandler {
	h := &PodSyncHandler{}

	h.w = w

	return h
}
