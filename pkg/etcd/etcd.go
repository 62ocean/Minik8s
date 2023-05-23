package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"k8s/pkg/api/pod"
	"time"
)

var cli *clientv3.Client
var err error

func EtcdInit(endpoint string) {
	fmt.Printf("init etcd\n")
	cli = GetEtcdClient(endpoint)
}
func GetEtcdClient(endpoint string) *clientv3.Client {
	if endpoint == "" {
		endpoint = "localhost:2379"
	}
	if cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{endpoint},
	}); err != nil {
		panic("connect to etcd failed: " + err.Error())
		return nil
	} else {
		fmt.Printf("connect to etcd successfully with endpoint: %s", endpoint)
		return cli
	}
}

func Put(key string, val string) bool {
	fmt.Printf("put\n")
	ctx, cancel := context.WithCancel(context.Background())
	_, err := cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		fmt.Println("put to etcd failed", err)
		return false
	}
	return true
}

func Del(key string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	_, err := cli.Delete(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("delete key %s failed\n", key)
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

func WatchPrefix(prefix string, vx *pod.VxlanDevice, localIp string, Cli *clientv3.Client) {
	watcher := clientv3.NewWatcher(cli)
	ctx, _ := context.WithCancel(context.TODO())
	watchChan := watcher.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for w := range watchChan {
		for _, event := range w.Events {
			if event.Type == mvccpb.PUT {
				node := pod.NodeNetwork{}
				// 新增节点
				if event.PrevKv == nil {
					_ = json.Unmarshal(event.Kv.Value, &node)
					if node.IpAddr != localIp {
						fmt.Printf("New node add to the flannel network\n")
						vx.AddNodeToNetwork(node.Subnet, node.Gateway, node.Docker0MacAddr, node.IpAddr)
					}
				} else {
					_ = json.Unmarshal(event.PrevKv.Value, &node)
					if node.IpAddr != localIp {
						fmt.Printf("Node update\n")
						vx.DelNodeFromNetwork(node.Subnet, node.Gateway, node.Docker0MacAddr, node.IpAddr)
						_ = json.Unmarshal(event.Kv.Value, &node)
						vx.AddNodeToNetwork(node.Subnet, node.Gateway, node.Docker0MacAddr, node.IpAddr)
					}
				}
			} else if event.Type == mvccpb.DELETE {
				node := pod.NodeNetwork{}
				_ = json.Unmarshal(event.PrevKv.Value, &node)
				if node.IpAddr != localIp {
					fmt.Printf("Delete node from flannel network\n")
					vx.DelNodeFromNetwork(node.Subnet, node.Gateway, node.Docker0MacAddr, node.IpAddr)
				}
			}
		}
	}
}

func EtcdDeinit() {
	if cli != nil {
		cli.Close()
	}
	fmt.Println("connect to etcd close")
}
