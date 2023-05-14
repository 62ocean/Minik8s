package listeners

import (
	storagepb2 "github.com/coreos/etcd/storage/storagepb"
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
func (s ServiceListener) OnSet(kv storagepb2.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

func (s ServiceListener) OnModify(kv storagepb2.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
}

func (p ServiceListener) OnDelete(kv storagepb2.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(kv.Key) + "\n")
}
