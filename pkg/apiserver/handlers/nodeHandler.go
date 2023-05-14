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

func CreateNode(request *restful.Request, response *restful.Response) {
	log.Println("Get create node request")
	node := new(object.Node)
	err := request.ReadEntity(&node)
	if err != nil {
		log.Println(err)
		return
	}
	ip := node.IP
	key := "/registry/nodes/default/" + ip
	nodeStorage := object.NodeStorage{
		Node:   *node,
		Status: object.RUNNING,
	}
	nodeString, _ := json.Marshal(nodeStorage)
	res := etcd.Put(key, string(nodeString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Node could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		ret := "ok"
		_, err := response.Write([]byte(ret))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetNode(request *restful.Request, response *restful.Response)    {}
func UpdateNode(request *restful.Request, response *restful.Response) {}
func RemoveNode(request *restful.Request, response *restful.Response) {}
