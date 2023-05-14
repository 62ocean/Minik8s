package listeners

import (
	"encoding/json"
	"fmt"
	storagepb2 "github.com/coreos/etcd/storage/storagepb"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	log "log"
)

/*-----------------Replicaset Etcd Listener---------------*/

type ReplicasetListener struct {
	publisher *publisher.Publisher
}

func NewReplicasetListener() *ReplicasetListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := ReplicasetListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Replicaset Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p ReplicasetListener) OnSet(kv storagepb2.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnModify etcd中对应资源被修改时回调
func (p ReplicasetListener) OnModify(kv storagepb2.KeyValue) {
	log.Printf("ETCD: modify key:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg, err := json.Marshal(kv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(kv)
	err = p.publisher.Publish("replicasets", jsonMsg, "PUT")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p ReplicasetListener) OnDelete(kv storagepb2.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(kv.Key) + "\n")
	jsonMsg, err := json.Marshal(kv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(kv)
	err = p.publisher.Publish("replicasets", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
