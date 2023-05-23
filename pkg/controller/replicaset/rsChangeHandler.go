package replicaset

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"log"
)

type RSChangeHandler struct {
	c *controller
}

func (h *RSChangeHandler) Handle(msg []byte) {

	log.Println("replicaset receive msg: " + string(msg))

	var msgObject object.MQMessage
	err := json.Unmarshal(msg, &msgObject)
	if err != nil {
		fmt.Println("[rscontroller] unmarshall msg failed")
		return
	}

	var rs object.ReplicaSet
	if msgObject.EventType == object.DELETE {
		err = json.Unmarshal([]byte(msgObject.PrevValue), &rs)
	} else {
		err = json.Unmarshal([]byte(msgObject.Value), &rs)
	}

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
