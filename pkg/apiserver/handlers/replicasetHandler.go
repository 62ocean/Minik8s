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

func GetReplicaset(request *restful.Request, response *restful.Response)    {}
func UpdateReplicaset(request *restful.Request, response *restful.Response) {}
func RemoveReplicaset(request *restful.Request, response *restful.Response) {}
