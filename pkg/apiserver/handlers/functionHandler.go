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

func CreateFunction(request *restful.Request, response *restful.Response) {
	log.Println("Get create function request")
	fun := new(object.Function)
	err := request.ReadEntity(&fun)
	if err != nil {
		log.Println(err)
		return
	}
	key := "/registry/functions/serverless/" + fun.Name
	funString, _ := json.Marshal(fun)
	res := etcd.Put(key, string(funString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Function could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		podQueue := "functions"
		//err := response.WriteEntity(string(podQueue))
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetAllFunction(request *restful.Request, response *restful.Response) {
	funMap := etcd.GetDirectory("/registry/functions")
	msg, _ := json.Marshal(funMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetFunction(request *restful.Request, response *restful.Response) {}
func UpdateFunction(request *restful.Request, response *restful.Response) {
	newFunInfo := object.Function{}
	err := request.ReadEntity(&newFunInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}
	newVal, _ := json.Marshal(&newFunInfo)
	key := "/registry/functions/serverless/" + newFunInfo.Name
	var ret string
	if etcd.GetOne(key) == "" {
		ret = "non-existed function"
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
func RemoveFunction(request *restful.Request, response *restful.Response) {}
