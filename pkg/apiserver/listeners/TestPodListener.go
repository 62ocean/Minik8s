package listeners

import (
	"fmt"
	storagepb2 "github.com/coreos/etcd/storage/storagepb"
	"k8s/pkg/etcd"
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
	etcd.EtcdInit("")
	etcd.Put("/registry/pod", "new pod 2 lalalalla")
	etcd.Put("/registry/node", "lalallalal")

}

func (p podListener) OnSet(kv storagepb2.KeyValue) {
	fmt.Printf("set watcher of key " + string(kv.Key) + "\n")
	return
}
func (p podListener) OnModify(kv storagepb2.KeyValue) {
	fmt.Printf("modify kye:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	return
}

func (p podListener) OnDelete(kv storagepb2.KeyValue) {
	fmt.Printf("delete kye:" + string(kv.Key) + "\n")
	return
}
