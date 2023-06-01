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

func CreateWorkflow(request *restful.Request, response *restful.Response) {
	log.Println("Get create workflow request")
	wf := new(object.Workflow)
	err := request.ReadEntity(&wf)
	if err != nil {
		log.Println(err)
		return
	}
	key := "/registry/workflows/serverless/" + wf.Metadata.Name
	funString, _ := json.Marshal(wf)
	res := etcd.Put(key, string(funString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Workflow could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		podQueue := "workflows"
		//err := response.WriteEntity(string(podQueue))
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetAllWorkflow(request *restful.Request, response *restful.Response) {
	funMap := etcd.GetDirectory("/registry/workflows")
	msg, _ := json.Marshal(funMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetWorkflow(request *restful.Request, response *restful.Response) {}
func UpdateWorkflow(request *restful.Request, response *restful.Response) {
	newWfInfo := object.Workflow{}
	err := request.ReadEntity(&newWfInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}
	newVal, _ := json.Marshal(&newWfInfo)
	key := "/registry/workflows/serverless/" + newWfInfo.Metadata.Name
	var ret string
	if etcd.GetOne(key) == "" {
		ret = "non-existed workflow"
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
func RemoveWorkflow(request *restful.Request, response *restful.Response) {}
