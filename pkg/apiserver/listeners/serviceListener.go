package listeners

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

/*-----------------Service Etcd Listener---------------*/

type ServiceListener struct {
	publisher *publisher.Publisher
}

func NewServiceListener() *ServiceListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := ServiceListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Service Etcd Handler-----------------*/

func (s ServiceListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

func (s ServiceListener) OnCreate(kv mvccpb.KeyValue) {
	fmt.Printf("ServiceListener: create key:" + string(kv.Key) + "value:" + string(kv.Value) + "\n")
	service := object.Service{}
	json.Unmarshal(kv.Value, &service)

	// 筛选所有标签符合的pod，构造对应该service的endpoint
	podMap := new(map[string]string)
	*podMap = etcd.GetDirectory("/registry/pods")
	var allPodsStorage []object.PodStorage
	fmt.Printf("list len %d\n", len(*podMap))
	for _, v := range *podMap {
		podStorage := object.PodStorage{}
		_ = json.Unmarshal([]byte(v), &podStorage)
		allPodsStorage = append(allPodsStorage, podStorage)
	}
	fmt.Println(len(allPodsStorage))

	endpoint := object.Endpoint{
		ServiceName: service.Metadata.Name,
		Selector:    service.Spec.Selector,
	}
	endpoint.PodIps = make(map[string]string)
	for _, podSto := range allPodsStorage {
		if service.Spec.Selector.App == podSto.Config.Metadata.Labels.App {
			if service.Spec.Selector.Env == podSto.Config.Metadata.Labels.Env {
				podIp := podSto.Config.IP
				podId := podSto.Config.Metadata.Uid
				fmt.Println(podId)
				fmt.Println(podIp)
				endpoint.PodIps[podId] = podIp
				// fmt.Println(endpoint.PodIps[podId])
			}
		}
	}
	fmt.Println(len(endpoint.PodIps))
	// 将endpoint注册进etcd

	endpointByte, _ := json.Marshal(endpoint)
	key := "/registry/endpoints/" + service.Metadata.Name
	etcd.Put(key, string(endpointByte))

	// 向“services”队列发布CREATE消息
	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	var err error
	err = s.publisher.Publish("services", jsonMsg, "CREATE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func (s ServiceListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete key:" + string(prevkv.Key) + " " + string(prevkv.Value) + "\n")
	log.Printf("ETCD: delete kv:" + string(kv.Key) + " " + string(kv.Value) + "\n")
	service := object.Service{}
	json.Unmarshal(prevkv.Value, &service)

	// 向“services”队列发布DELETE消息
	jsonMsg := publisher.ConstructPublishMsg(prevkv, prevkv, object.DELETE)
	var err error
	err = s.publisher.Publish("services", jsonMsg, "DELETE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 在etcd中删除service对应的endpoint
	key := "/registry/endpoints/" + service.Metadata.Name
	res := etcd.Del(key)
	if !res {
		fmt.Println("serviceOnDelete: delete endpoint failed!")
	}

	return
}

func (s ServiceListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modify key:" + string(prevkv.Key) + " value:" + string(prevkv.Value) + "\n")

	etcd.Del(string(prevkv.Key))
	etcd.Put(string(kv.Key), string(kv.Value))

}
