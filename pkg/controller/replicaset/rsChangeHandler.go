package replicaset

import (
	"encoding/json"
	"fmt"
	"k8s/object"
)

type RSChangeHandler struct {
	c *controller
}

func (h *RSChangeHandler) Handle(msg []byte) {

	fmt.Println("replicaset receive msg: " + string(msg))

	var msgObject object.MQMessage
	err := json.Unmarshal(msg, &msgObject)
	if err != nil {
		fmt.Println("[rscontroller] unmarshall msg failed")
		return
	}

	fmt.Println(msgObject)

	var rs object.ReplicaSet
	err = json.Unmarshal([]byte(msgObject.Value), &rs)
	if err != nil {
		fmt.Println("[rscontroller] unmarshall changed replicaset failed")
		return
	}

	//var msg_type int
	//var msg_rs object.ReplicaSet

	switch msgObject.EventType {
	case object.CREATE:
		h.c.AddReplicaset(rs)
	case object.DELETE:
		h.c.DeleteReplicaset(rs)
	case object.UPDATE:
		h.c.UpdateReplicaset(rs)
	}
}

func NewReplicasetChangeHandler(c *controller) *RSChangeHandler {
	h := &RSChangeHandler{}

	h.c = c

	return h
}
