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

func CreatePod(request *restful.Request, response *restful.Response) {
	log.Println("Get create pod request")
	pod := new(object.Pod)
	err := request.ReadEntity(&pod)
	if err != nil {
		log.Println(err)
		return
	}
	uid := pod.Metadata.Uid
	key := "/registry/pods/default/" + uid
	podStorage := object.PodStorage{
		Config: *pod,
		Status: object.STOPPED,
	}
	podString, _ := json.Marshal(podStorage)
	res := etcd.Put(key, string(podString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Pod could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		podQueue := "pods"
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetPod(request *restful.Request, response *restful.Response) {}
func UpdatePod(request *restful.Request, response *restful.Response) {
	newPodInfo := object.PodStorage{}
	err := request.ReadEntity(&newPodInfo)
	fmt.Println(newPodInfo)
	if err != nil {
		log.Println(err)
		return
	}
	newVal, _ := json.Marshal(&newPodInfo)
	key := "/registry/pods/default/" + newPodInfo.Config.Metadata.Uid
	var ret string
	if etcd.GetOne(key) == "" {
		ret = "non-existed pod"
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
func RemovePod(request *restful.Request, response *restful.Response) {}
func GetAllPod(request *restful.Request, response *restful.Response) {
	podMap := etcd.GetDirectory("/registry/pods")
	msg, _ := json.Marshal(podMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}
