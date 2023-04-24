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

func (p podListener) Set([]byte, []byte) {
	return
}

func (p podListener) Create([]byte, []byte) {
	return
}
func (p podListener) Modify(a []byte, b []byte) {
	fmt.Printf("modify kye:" + string(a) + "value:" + string(b) + "\n")
	return
}

func (p podListener) Delete(a []byte) {
	fmt.Printf("delete kye:" + string(a) + "\n")
	return
}
