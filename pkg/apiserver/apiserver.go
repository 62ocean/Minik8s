package apiserver

import (
	"fmt"
	"k8s/pkg/etcd"
)

func StartServer() {
	//etcd.EtcdTest()
	etcd.EtcdInit()
}

func CreatePod() {
	fmt.Printf("apiserver: create pod\n")
}

func DeletePod() {
	fmt.Printf("apiserver: delete pod\n")
}

func DescribePod() {
	fmt.Printf("apiserver: describe pod\n")
}

func DescribeService() {
	fmt.Printf("apiserver: describe service\n")
}

func EtcdGetOne(key string) string {
	res := etcd.GetOne(key)
	if res == "" {
		fmt.Printf("get key %s from etcd failed", key)
	}
	return res
}

func EtcdGetDirectory(prefix string) map[string]string {
	res := etcd.GetDirectory(prefix)
	if res == nil {
		fmt.Printf("get directory %s from etcd failed", prefix)
	}
	return res
}

func EtcdPut(key string, val string) bool {
	res := etcd.Put(key, val)
	return res
}
