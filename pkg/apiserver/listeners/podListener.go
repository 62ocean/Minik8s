package listeners

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	log "log"
)

/*-----------------Pod Etcd Listener---------------*/

type PodListener struct {
	publisher *publisher.Publisher
}

func NewPodListener() *PodListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := PodListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Pod Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p PodListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnCreate etcd中对应资资源被创建时回调
func (p PodListener) OnCreate(kv mvccpb.KeyValue) {
	log.Printf("ETCD: create kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	podStorage := object.PodStorage{}
	_ = json.Unmarshal(kv.Value, &podStorage)
	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	var err error
	// forward to relicaset
	log.Println("publish CREATE to pods")
	err = p.publisher.Publish("pods", jsonMsg, "CREATE")
	// forward to kubelet
	if podStorage.Node != "" {
		log.Println("publish CREATE to pods_node")
		err = p.publisher.Publish("pods_node", jsonMsg, "CREATE")
	}
	// forward to scheduler
	if podStorage.Node == "" {
		log.Println("publish CREATE to pods_sched")
		err = p.publisher.Publish("pods_sched", jsonMsg, "CREATE")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnModify etcd中对应资源被修改时回调
func (p PodListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	podStorage := object.PodStorage{}
	_ = json.Unmarshal(kv.Value, &podStorage)
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.UPDATE)
	var err error
	// forward to relicaset
	log.Println("publish PUT to pods")
	err = p.publisher.Publish("pods", jsonMsg, "PUT")
	// forward to kubelet
	if podStorage.Node != "" {
		log.Println("publish PUT to pods_node")
		err = p.publisher.Publish("pods_node", jsonMsg, "PUT")
	}
	// forward to scheduler
	if podStorage.Node == "" {
		log.Println("publish PUT to pods_sched")
		err = p.publisher.Publish("pods_sched", jsonMsg, "PUT")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p PodListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.DELETE)
	var err error
	// forward to relicaset
	log.Println("publish DEL to pods")
	err = p.publisher.Publish("pods", jsonMsg, "DEL")
	// forward to kubelet
	log.Println("publish DEL to pods_node")
	err = p.publisher.Publish("pods_node", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
