package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreatePod(request *restful.Request, response *restful.Response) {
	log.Println("Get create pod request")
	pod := new(object.Pod)
	err := request.ReadEntity(&pod)
	if err != nil {
		log.Println(err)
		return
	}
	name := pod.Metadata.Name
	replica := getReplicaIndex(name)
	podStorage := object.PodStorage{
		Config:  *pod,
		Status:  object.STOPPED,
		Replica: replica,
	}
	key := "/registry/pods/default/" + name + "-" + strconv.Itoa(replica)
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
		//err := response.WriteEntity(string(podQueue))
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetPod(request *restful.Request, response *restful.Response) {}
func UpdatePod(request *restful.Request, response *restful.Response) {
	//newPodInfo := object.PodStorage{}
	//err := request.ReadEntity(newPodInfo)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//newVal, _ := json.Marshal(newPodInfo)
	//key := "/registry/pods/" + newPodInfo
	//etcd.Put()

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

func getReplicaIndex(name string) int {
	lastHyphen := strings.LastIndex(name, "-")
	replica, err := strconv.Atoi(name[lastHyphen+1:])
	if err != nil || replica == -1 {
		replica = 0
	}
	log.Println("Get index from name: " + name + ", replica: " + strconv.Itoa(replica))
	return replica
}
