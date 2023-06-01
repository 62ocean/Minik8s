package handlers

import (
	"encoding/json"
	"fmt"
	"k8s/pkg/etcd"
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func GetEndpoint(request *restful.Request, response *restful.Response) {
	var serviceName string
	request.ReadEntity(&serviceName)
	fmt.Printf("apiserver handler: get endpoint %s\n", serviceName)

	serviceStr := etcd.GetOne("/registry/endpoints/" + serviceName)
	msg, _ := json.Marshal(serviceStr)
	_, err := response.Write(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func RemoveEndpoint(request *restful.Request, response *restful.Response) {
	log.Printf("apiserver handler: delete endpoint")
	var serviceName string
	request.ReadEntity(&serviceName)
	key := "/registry/endpoints/" + serviceName
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
