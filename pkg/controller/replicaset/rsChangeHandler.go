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

	var rs object.ReplicaSet
	err = json.Unmarshal([]byte(msgObject.Value), &rs)
	if err != nil {
		fmt.Println("[rscontroller] unmarshall changed replicaset failed")
		return
	}

	h.c.ReplicasetChangeHandler(msgObject.EventType, rs)
}

func NewReplicasetChangeHandler(c *controller) *RSChangeHandler {
	h := &RSChangeHandler{}
	h.c = c

	return h
}
