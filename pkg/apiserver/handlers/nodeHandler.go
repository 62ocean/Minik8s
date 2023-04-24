package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/etcd"
	"net/http"
)

func RegisterNode(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/node").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.POST("/create").To(CreateNode))
	ws.Route(ws.GET("/get").To(GetNode))
	ws.Route(ws.POST("/update").To(UpdateNode))
	ws.Route(ws.DELETE("/remove").To(RemoveNode))
	container.Add(ws)
}

func CreateNode(request *restful.Request, response *restful.Response) {
	fmt.Println(request.Request.Host)
	fmt.Println(request.Request.Header)
	fmt.Println(request.Request.Form)
	// TODO: 获取请求发送者的ip
	ip := "127.0.0.1"
	key := "/registry/minions/" + ip
	val := request.PathParameter("info")
	res := etcd.Put(key, val)
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		response.WriteErrorString(http.StatusNotFound, "node could not be persisted")
	} else {
		response.WriteEntity("ok")
	}
}

func GetNode(request *restful.Request, response *restful.Response)    {}
func UpdateNode(request *restful.Request, response *restful.Response) {}
func RemoveNode(request *restful.Request, response *restful.Response) {}
