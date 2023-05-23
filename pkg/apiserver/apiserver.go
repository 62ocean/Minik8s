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
	serviceListener    *listeners.ServiceListener
	nodeListener       *listeners.NodeListener
	hpaListener        *listeners.HpaListener

	//TODO 在此添加其他listener……
}

// CreateAPIServer 初始化APIServer结构体中的内容
func CreateAPIServer() (*APIServer, error) {
	// etcd watcher
	etcd.EtcdInit(global.EtcdHost)
	etcdWatcher, err := etcd.NewEtcdWatcher([]string{global.EtcdHost})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}

	// listeners
	podListener := listeners.NewPodListener()
	replicasetListener := listeners.NewReplicasetListener()
	serviceListener := listeners.NewServiceListener()
	nodeListener := listeners.NewNodeListener()
	hpaListener := listeners.NewHpaListener()

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
		serviceListener:    serviceListener,
		nodeListener:       nodeListener,
		hpaListener:        hpaListener,
	}

	return &server, nil
}

// StartServer 开始监听相关端口请求
func (s *APIServer) StartServer() {
	// watch
	s.etcdWatcher.AddWatch("/registry/pods/", true, s.podListener)
	s.etcdWatcher.AddWatch("/registry/replicasets/", true, s.replicasetListener)
	s.etcdWatcher.AddWatch("/registry/services/", true, s.serviceListener)
	s.etcdWatcher.AddWatch("/registry/nodes/", true, s.nodeListener)
	s.etcdWatcher.AddWatch("/registry/hpas/", true, s.hpaListener)

	// list
	server := &http.Server{Addr: ":8080", Handler: s.wsContainer}
	defer server.Close()
	log.Fatal(server.ListenAndServe())
}
