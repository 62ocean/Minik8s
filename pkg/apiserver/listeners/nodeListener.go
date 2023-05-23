package listeners

import (
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	"log"
)

/*-----------------Pod Etcd Listener---------------*/

type NodeListener struct {
	publisher *publisher.Publisher
}

func NewNodeListener() *NodeListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := NodeListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Pod Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p NodeListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnCreate etcd中对应资资源被创建时回调
func (p NodeListener) OnCreate(kv mvccpb.KeyValue) {
	log.Printf("ETCD: create kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	// forward to scheduler
	err := p.publisher.Publish("nodes", jsonMsg, "CREATE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnModify etcd中对应资源被修改时回调
func (p NodeListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.UPDATE)
	// forward to scheduler
	err := p.publisher.Publish("nodes", jsonMsg, "PUT")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p NodeListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.DELETE)
	// forward to scheduler
	err := p.publisher.Publish("nodes", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
