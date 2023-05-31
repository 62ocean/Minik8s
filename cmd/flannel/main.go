package main

import (
	"k8s/pkg/apiserver/flannel"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
)

func main() {
	etcd.EtcdInit(global.EtcdHost)
	flannel.Exec()
}
