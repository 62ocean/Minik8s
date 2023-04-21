package main

import (
	"k8s/pkg/apiserver"
	"k8s/pkg/kubectl"
)

func main() {
	//etcd.EtcdTest()
	apiserver.StartServer()
	kubectl.CmdExec()
}
