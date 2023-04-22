package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var cli *clientv3.Client
var err error

func EtcdInit() {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect to etcd failed: ", err)
		return
	}
	fmt.Println("connect to etcd success")
}

func Put(key string, val string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	_, err := cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		fmt.Println("put to etcd failed", err)
		return false
	}
	return true
}

// 获取key对应的value
func GetOne(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed: ", err)
		return ""
	}
	length := len(resp.Kvs)
	if length != 1 {
		fmt.Printf("with result length %+v which should be 1", length)
		return ""
	}
	return string(resp.Kvs[0].Value)
}

// 获取key的前缀为prefix的所有键值对
func GetDirectory(prefix string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	getResp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed: ", err)
		return nil
	}
	res := make(map[string]string)
	for _, resp := range getResp.Kvs {
		//fmt.Printf("key %s, value:%s\n", string(resp.Key), string(resp.Value))
		res[string(resp.Key)] = string(resp.Value)
	}
	return res
}

func EtcdDeinit() {
	cli.Close()
	fmt.Println("connect to etcd close")
}
