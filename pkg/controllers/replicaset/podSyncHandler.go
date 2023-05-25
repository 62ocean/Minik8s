package replicaset

import "log"

type PodSyncHandler struct {
	w *worker
}

func (h *PodSyncHandler) Handle(msg []byte) {

	log.Println("[rs worker] receive pod change msg")

	//var msgObject object.MQMessage
	//err := json.Unmarshal(msg, &msgObject)
	//if err != nil {
	//	fmt.Println("[worker] unmarshall msg failed")
	//	return
	//}
	//
	//var podStorage object.PodStorage
	//err = json.Unmarshal([]byte(msgObject.Value), &podStorage)
	//if msgObject.EventType == object.DELETE {
	//	err = json.Unmarshal([]byte(msgObject.PrevValue), &podStorage)
	//} else {
	//	err = json.Unmarshal([]byte(msgObject.Value), &podStorage)
	//}
	//if err != nil {
	//	fmt.Println("[worker] unmarshall changed pod failed")
	//	return
	//}

	//if podStorage.Config.Metadata.Labels.App == h.w.target.Spec.Selector.MatchLabels.App &&
	//	podStorage.Config.Metadata.Labels.Env == h.w.target.Spec.Selector.MatchLabels.Env {
	h.w.SyncPods()
	//}
}

func NewPodSyncHandler(w *worker) *PodSyncHandler {
	h := &PodSyncHandler{}

	h.w = w

	return h
}
