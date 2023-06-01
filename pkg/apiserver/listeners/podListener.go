package listeners

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	log "log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

/*-----------------Pod Etcd Listener---------------*/

type PodListener struct {
	publisher *publisher.Publisher
}

func NewPodListener() *PodListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := PodListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Pod Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p PodListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnCreate etcd中对应资资源被创建时回调
func (p PodListener) OnCreate(kv mvccpb.KeyValue) {
	log.Printf("ETCD: create kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	podStorage := object.PodStorage{}

	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	var err error

	// forward to relicaset
	log.Println("publish CREATE to pods_XXX")
	_ = json.Unmarshal(kv.Value, &podStorage)
	err = p.publisher.Publish("pods_"+podStorage.Config.Metadata.Labels.App, jsonMsg, "CREATE")

	//_ = json.Unmarshal(kv.Value, &podStorage)
	// forward to kubelet
	if podStorage.Node != "" {
		log.Println("publish CREATE to pods_node")
		err = p.publisher.Publish("pods_node", jsonMsg, "CREATE")
	}
	// forward to scheduler
	if podStorage.Node == "" {
		log.Println("publish CREATE to pods_sched")
		err = p.publisher.Publish("pods_sched", jsonMsg, "CREATE")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// epMap := etcd.GetDirectory("/registry/endpoints")
	// for k, v := range epMap {
	// 	ep := object.Endpoint{}
	// 	json.Unmarshal([]byte(v), &ep)
	// 	if (podStorage.Config.Metadata.Labels.App == ep.Selector.App) &&
	// 		(podStorage.Config.Metadata.Labels.Env == ep.Selector.Env) {
	// 		id := podStorage.Config.Metadata.Uid
	// 		ip := podStorage.Config.IP
	// 		ep.PodIps[id] = ip
	// 		epByte, _ := json.Marshal(ep)
	// 		etcd.Put(k, string(epByte))
	// 	}
	// }

	//err = p.publisher.Publish("pods", jsonMsg, "CREATE")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	return
}

// OnModify etcd中对应资源被修改时回调
func (p PodListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	fmt.Printf("ETCD: modify kye:" + string(kv.Key) + "\nprevvalue:" + string(prevkv.Value) + "\n" + "\nvalue:" + string(kv.Value) + "\n")
	podStorage := object.PodStorage{}
	prevPodStorage := object.PodStorage{}
	_ = json.Unmarshal(kv.Value, &podStorage)
	_ = json.Unmarshal(prevkv.Value, &prevPodStorage)
	if podStorage.Node == prevPodStorage.Node && podStorage.Config.IP == prevPodStorage.Config.IP {
		return
	}

	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.UPDATE)
	var err error
	// forward to relicaset
	log.Println("publish PUT to pods_XXX")
	exchangeName1 := "pods_" + podStorage.Config.Metadata.Labels.App
	err = p.publisher.Publish(exchangeName1, jsonMsg, "PUT")
	exchangeName2 := "pods_" + prevPodStorage.Config.Metadata.Labels.App
	if exchangeName1 != exchangeName2 {
		err = p.publisher.Publish(exchangeName2, jsonMsg, "PUT")
	}
	// forward to kubelet
	if podStorage.Node != "" {
		log.Println("publish UPDATE to pods_node")
		err = p.publisher.Publish("pods_node", jsonMsg, "PUT")
	}
	// forward to scheduler
	if podStorage.Node == "" {
		log.Println("publish UPDATE to pods_sched")
		err = p.publisher.Publish("pods_sched", jsonMsg, "PUT")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// pod 完成调度
	fmt.Printf("Pod Listener: prevIP %s\n", prevPodStorage.Config.IP)
	fmt.Printf("Pod Listener: curIP %s\n", podStorage.Config.IP)
	if podStorage.Config.IP != "" && podStorage.Config.IP != prevPodStorage.Config.IP {
		fmt.Println("add endpoint")
		// 遍历endpoint，向符合label的endpoint插入新增pod的ip记录
		// 将新endpoint写回etcd
		epMap := etcd.GetDirectory("/registry/endpoints")
		for k, v := range epMap {
			ep := object.Endpoint{}
			json.Unmarshal([]byte(v), &ep)
			if (podStorage.Config.Metadata.Labels.App == ep.Selector.App) &&
				(podStorage.Config.Metadata.Labels.Env == ep.Selector.Env) {
				id := podStorage.Config.Metadata.Uid
				ip := podStorage.Config.IP
				ep.PodIps[id] = ip
				epByte, _ := json.Marshal(ep)
				etcd.Put(k, string(epByte))
			}
		}

	}

	//err = p.publisher.Publish("pods", jsonMsg, "PUT")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p PodListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.DELETE)
	var err error
	var podStorage object.PodStorage

	// forward to relicaset
	log.Println("publish DEL to pods_XXX")
	_ = json.Unmarshal([]byte(prevkv.Value), &podStorage)
	err = p.publisher.Publish("pods_"+podStorage.Config.Metadata.Labels.App, jsonMsg, "DEL")

	// forward to kubelet
	log.Println("publish DEL to pods_node")
	err = p.publisher.Publish("pods_node", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 把该pod从相关service的endpoint里删掉
	// 此处podStorage存的是删掉的值
	epMap := etcd.GetDirectory("/registry/endpoints")
	for k, v := range epMap {
		ep := object.Endpoint{}
		json.Unmarshal([]byte(v), &ep)
		if (podStorage.Config.Metadata.Labels.App == ep.Selector.App) &&
			(podStorage.Config.Metadata.Labels.Env == ep.Selector.Env) {
			id := podStorage.Config.Metadata.Uid
			delete(ep.PodIps, id)
			epByte, _ := json.Marshal(ep)
			etcd.Put(k, string(epByte))
		}
	}

	//err = p.publisher.Publish("pods", jsonMsg, "DEL")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	return
}
