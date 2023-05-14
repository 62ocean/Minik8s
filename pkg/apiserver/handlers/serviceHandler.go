package handlers

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
)

func CreateService(request *restful.Request, response *restful.Response) {
	log.Printf("apiserver handler: create service")

	service := new(object.Service)
	err := request.ReadEntity(&service)

	if err != nil {
		log.Println(err)
		return
	}

	serviceName := service.Metadata.Name
	key := "/registry/services/" + serviceName
	serviceString, _ := json.Marshal(*service)
	res := etcd.Put(key, string(serviceString))
	response.AddHeader("Content-Type", "text/plain")
	
}

func GetService(request *restful.Request, response *restful.Response)    {}
func UpdateService(request *restful.Request, response *restful.Response) {}
func RemoveService(request *restful.Request, response *restful.Response) {}
