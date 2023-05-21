package listeners

import (
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
)

// 实现etcd watcher listener的监控处理函数
type podListener struct {
}

func main() {
	test := podListener{}
	listener := etcd.Listener(test)
	watch, _ := etcd.NewEtcdWatcher([]string{"localhost:2379"})
	watch.AddWatch("/registry", true, listener)
	defer watch.Close(true)
	etcd.EtcdInit(global.EtcdHost)
	etcd.Put("/registry/pod", "new pod 2 lalalalla")
	etcd.Put("/registry/node", "lalallalal")

}

func (p podListener) OnSet(kv mvccpb.KeyValue) {
	fmt.Printf("set watcher of key " + string(kv.Key) + "\n")
	return
}
func (p podListener) OnCreate(kv mvccpb.KeyValue) {
	fmt.Printf("create kye:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	return
}

func (p podListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("modify kye:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	return
}

func (p podListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	fmt.Printf("delete kye:" + string(kv.Key) + "\n")
	return
}
