package listeners

import (
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DnsListener struct {
	publisher *publisher.Publisher
}

func NewDnsListener() *DnsListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := DnsListener{
		publisher: newPublisher,
	}
	return &listener
}

func (d DnsListener) OnSet(kv mvccpb.KeyValue) {
	fmt.Printf("set watcher of key " + string(kv.Key) + "\n")
	return
}
func (d DnsListener) OnCreate(kv mvccpb.KeyValue) {
	fmt.Printf("create key:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")

	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	err := d.publisher.Publish("dns", jsonMsg, "CREATE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func (d DnsListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("modify key:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	return
}

func (d DnsListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("delete key:" + string(kv.Key) + "\n")
	return
}
