package handlers

import (
	"fmt"
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
	etcd.EtcdInit()
	etcd.Put("/registry/pod", "new pod 2 lalalalla")
	etcd.Put("/registry/node", "lalallalal")

}

func (p podListener) OnSet(key []byte, val []byte) {
	fmt.Printf("set watcher of key " + string(key) + "\n")
	return
}
func (p podListener) OnModify(key []byte, val []byte) {
	fmt.Printf("modify kye:" + string(key) + "value:" + string(val) + "\n")
	return
}

func (p podListener) OnDelete(key []byte) {
	fmt.Printf("delete kye:" + string(key) + "\n")
	return
}
