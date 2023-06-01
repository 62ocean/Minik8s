package handlers

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
)

var cnt = 0

func getFreeClusterIP() string {
	cnt++
	return fmt.Sprintf("10.111.111.%d", cnt)

}

// 请求参数为service实体
func CreateService(request *restful.Request, response *restful.Response) {
	log.Printf("apiserver handler: create service")

	service := new(object.Service)
	err := request.ReadEntity(&service)
	if err != nil {
		log.Println(err)
		return
	}
	id, _ := uuid.NewUUID()
	service.Metadata.Uid = id.String()
	fmt.Printf("service id : %s\n", service.Metadata.Uid)

	if service.Spec.ClusterIP == "" {
		service.Spec.ClusterIP = getFreeClusterIP()
	}

	serviceName := service.Metadata.Name
	key := "/registry/services/" + serviceName
	serviceByte, _ := json.Marshal(*service)
	res := etcd.Put(key, string(serviceByte))
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
	//service2.CreateService(*service)
}

func GetService(request *restful.Request, response *restful.Response) {
	var serviceName string
	request.ReadEntity(&serviceName)
	log.Printf("apiserver handler: get service %s\n", serviceName)

	serviceStr := etcd.GetOne("/registry/services/" + serviceName)
	msg, _ := json.Marshal(serviceStr)
	_, err := response.Write(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func UpdateService(request *restful.Request, response *restful.Response) {}

func GetAllService(request *restful.Request, response *restful.Response) {
	serviceMap := etcd.GetDirectory("/registry/services")
	msg, _ := json.Marshal(serviceMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 请求参数为service name
func RemoveService(request *restful.Request, response *restful.Response) {
	log.Printf("apiserver handler: delete service")
	var serviceName string
	request.ReadEntity(&serviceName)
	key := "/registry/services/" + serviceName
	res := etcd.Del(key)
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "Service could not be delete")
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

}

func CheckService(request *restful.Request, response *restful.Response) {
	serviceName := request.PathParameter("serviceName")

	endpointStr := etcd.GetOne("/registry/endpoints/" + serviceName)
	endpoint := object.Endpoint{}
	json.Unmarshal([]byte(endpointStr), &endpoint)
	flag := 0
	for _, v := range endpoint.PodIps {
		if v != "" {
			flag = 1
		}
	}

	msg, _ := json.Marshal(flag)
	_, err := response.Write(msg)
	if err != nil {
		fmt.Println(err.Error())
	}

}
