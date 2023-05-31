package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/etcd"
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
