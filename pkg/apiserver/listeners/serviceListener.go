package listeners

import (
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	"log"
)

/*-----------------Service Etcd Listener---------------*/

type ServiceListener struct {
	publisher *publisher.Publisher
}

func NewServiceListener() *ServiceListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := ServiceListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Service Etcd Handler-----------------*/

func (s ServiceListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

func (s ServiceListener) OnCreate(kv mvccpb.KeyValue) {
	fmt.Printf("create kye:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	return
}

func (s ServiceListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(prevkv.Key) + " value:" + string(prevkv.Value) + "\n")
}

func (p ServiceListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
}
