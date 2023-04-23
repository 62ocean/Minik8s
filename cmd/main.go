package main

import (
	"k8s/pkg/etcd"
)

func main() {
	etcd.EtcdTest()
	//apiserver.StartServer()
	//kubectl.CmdExec()
}
