package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"k8s/pkg/etcd"
	service2 "k8s/pkg/service"
	"log"
	"net/http"
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

	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Service could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		serviceQueue := "services"
		_, err := response.Write([]byte(serviceQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	service2.CreateService(*service)
}

func GetService(request *restful.Request, response *restful.Response)    {}
func UpdateService(request *restful.Request, response *restful.Response) {}
func RemoveService(request *restful.Request, response *restful.Response) {}
