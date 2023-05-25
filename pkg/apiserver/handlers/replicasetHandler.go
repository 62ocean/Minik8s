package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
	"net/http"
)

func CreateReplicaset(request *restful.Request, response *restful.Response) {
	log.Println("Get create replicaset request")
	rs := new(object.ReplicaSet)
	err := request.ReadEntity(&rs)
	if err != nil {
		log.Println(err)
		return
	}
	uid := rs.Metadata.Uid
	key := "/registry/replicasets/default/" + uid
	rsString, _ := json.Marshal(rs)
	res := etcd.Put(key, string(rsString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Replicaset could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		podQueue := "replicasets"
		//err := response.WriteEntity(string(podQueue))
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetAllReplicaset(request *restful.Request, response *restful.Response) {
	rsMap := etcd.GetDirectory("/registry/replicasets")
	msg, _ := json.Marshal(rsMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetReplicaset(request *restful.Request, response *restful.Response) {}
func UpdateReplicaset(request *restful.Request, response *restful.Response) {
	newRSInfo := object.ReplicaSet{}
	err := request.ReadEntity(&newRSInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}
	newVal, _ := json.Marshal(&newRSInfo)
	key := "/registry/replicasets/default/" + newRSInfo.Metadata.Uid
	var ret string
	if etcd.GetOne(key) == "" {
		ret = "non-existed rs"
		err1 := response.WriteErrorString(500, ret)
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	} else {
		etcd.Put(key, string(newVal))
		ret = "ok"
		_, err1 := response.Write([]byte(ret))
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	}
}
func RemoveReplicaset(request *restful.Request, response *restful.Response) {}
