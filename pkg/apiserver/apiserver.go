package apiserver

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/apiserver/handlers"
	"k8s/pkg/apiserver/listeners"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"log"
	"net/http"
)

/*-----------------APIServer-----------------*/

type APIServer struct {
	wsContainer        *restful.Container
	etcdWatcher        *etcd.EtcdWatcher
	podListener        *listeners.PodListener
	replicasetListener *listeners.ReplicasetListener
	//TODO 在此添加其他listener……
}

// CreateAPIServer 初始化APIServer结构体中的内容
func CreateAPIServer() (*APIServer, error) {
	// etcd watcher
	etcd.EtcdInit("")
	etcdWatcher, err := etcd.NewEtcdWatcher([]string{global.EtcdHost})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}

	// listeners
	podListener := listeners.NewPodListener()
	replicasetListener := listeners.NewReplicasetListener()

	// HTTP server
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	handlers.InitWebServer(wsContainer)

	// construct APIServer
	server := APIServer{
		etcdWatcher:        etcdWatcher,
		podListener:        podListener,
		wsContainer:        wsContainer,
		replicasetListener: replicasetListener,
	}

	return &server, nil
}

// StartServer 开始监听相关端口请求
func (s *APIServer) StartServer() {
	// watch
	s.etcdWatcher.AddWatch("/registry/pods/", true, s.podListener)
	s.etcdWatcher.AddWatch("/registry/replicasets/", true, s.replicasetListener)

	// list
	server := &http.Server{Addr: ":8080", Handler: s.wsContainer}
	defer server.Close()
	log.Fatal(server.ListenAndServe())
}

//
//func CreatePod() {
//	fmt.Printf("apiserver: create pod\n")
//}
//
//func DeletePod() {
//	fmt.Printf("apiserver: delete pod\n")
//}
//
//func DescribePod() {
//	fmt.Printf("apiserver: describe pod\n")
//}
//
//func DescribeService() {
//	fmt.Printf("apiserver: describe service\n")
//}
//
//func EtcdGetOne(key string) string {
//	res := etcd.GetOne(key)
//	if res == "" {
//		fmt.Printf("get key %s from etcd failed", key)
//	}
//	return res
//}
//
//func EtcdGetDirectory(prefix string) map[string]string {
//	res := etcd.GetDirectory(prefix)
//	if res == nil {
//		fmt.Printf("get directory %s from etcd failed", prefix)
//	}
//	return res
//}
//
//func EtcdPut(key string, val string) bool {
//	res := etcd.Put(key, val)
//	return res
//}
