package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/etcd"
	"log"
	"net/http"
)

func CreateNode(request *restful.Request, response *restful.Response) {
	log.Println("Get create node request")
	//ip := "127.0.0.1"
	//key := "/registry/nodes/" + ip
	//node := new(object.NodeStorage)
	//request.ReadEntity(&node)
	//nodestring, _ := json.Marshal(*node)
	//fmt.Println(string(nodestring))
	//res := etcd.Put(key, string(nodestring))
	key := "/registry/nodes/127.0.0.1"
	val := "{\"Node\":{\"Name\":\"TestNode\",\"IP\":\"127.0.0.1\"},\"Status\":0}"
	res := etcd.Put(key, val)
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Node could not be persisted")
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

func GetNode(request *restful.Request, response *restful.Response)    {}
func UpdateNode(request *restful.Request, response *restful.Response) {}
func RemoveNode(request *restful.Request, response *restful.Response) {}
