package main

import (
	"fmt"
	"k8s/pkg/etcd"
)

func main() {
	etcd.EtcdInit()
	etcd.Put("test", "test")
	etcd.Put("test1/test1_1", "test1_1")
	etcd.Put("test1/test1_2", "test1_2")
	etcd.Put("test2/test2_1", "test2_1")
	etcd.Put("test2/test2_2", "test2_2")
	etcd.Put("test2/test2_3", "test2_3")

	val := etcd.GetOne("test")
	fmt.Printf("test: %s\n", val)

	fmt.Printf("test1:\n")
	var slice map[string]string
	slice = etcd.GetDirectory("test1/")
	for k, v := range slice {
		fmt.Printf("%s: %s\n", k, v)
	}

	fmt.Printf("test2:\n")
	slice = etcd.GetDirectory("test2/")
	for k, v := range slice {
		fmt.Printf("%s: %s\n", k, v)
	}
	etcd.EtcdDeinit()
}
