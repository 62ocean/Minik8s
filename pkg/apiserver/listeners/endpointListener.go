package listeners

import (
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
)

type EndpointListener struct {
	publisher *publisher.Publisher
}

func NewEndpointListener() *EndpointListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := EndpointListener{
		publisher: newPublisher,
	}
	return &listener
}

func (e EndpointListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("ETCD: modify key:" + string(prevkv.Key) + " value:" + string(prevkv.Value) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.UPDATE)
	var err error
	err = e.publisher.Publish("endpoints", jsonMsg, "UPDATE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (e EndpointListener) OnSet(kv mvccpb.KeyValue) {
	fmt.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

func (e EndpointListener) OnCreate(kv mvccpb.KeyValue) {
	fmt.Printf("create key:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
}

func (e EndpointListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
}
