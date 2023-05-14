package handlers

import (
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
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

}

func GetService(request *restful.Request, response *restful.Response)    {}
func UpdateService(request *restful.Request, response *restful.Response) {}
func RemoveService(request *restful.Request, response *restful.Response) {}
