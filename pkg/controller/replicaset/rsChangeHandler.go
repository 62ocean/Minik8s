package replicaset

import "fmt"

type RSChangeHandler struct {
	c *controller
}

func (h *RSChangeHandler) Handle(msg []byte) {

	fmt.Println("replicaset msg: " + string(msg))

	//var msg_type int
	//var msg_rs object.ReplicaSet

	//switch msg_type {
	//case RS_CREATE:
	//	h.c.AddReplicaset(msg_rs)
	//case RS_DELETE:
	//	h.c.DeleteReplicaset(msg_rs)
	//case RS_UPDATE:
	//	h.c.UpdateReplicaset(msg_rs)
	//}
}

func NewReplicasetChangeHandler(c *controller) *RSChangeHandler {
	h := &RSChangeHandler{}

	h.c = c

	return h
}
