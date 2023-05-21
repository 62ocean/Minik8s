package handlers

import "github.com/emicklei/go-restful/v3"

func InitWebServer(container *restful.Container) {
	// node
	nodeWS := new(restful.WebService)
	nodeWS.Path("/nodes").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	nodeWS.Route(nodeWS.POST("/create").To(CreateNode))
	nodeWS.Route(nodeWS.GET("/get").To(GetNode))
	nodeWS.Route(nodeWS.POST("/update").To(UpdateNode))
	nodeWS.Route(nodeWS.DELETE("/remove").To(RemoveNode))
	container.Add(nodeWS)

	// pod
	podWS := new(restful.WebService)
	podWS.Path("/pods").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	podWS.Route(podWS.POST("/create").To(CreatePod))
	podWS.Route(podWS.GET("/get").To(GetPod))
	//podWS.Route(podWS.POST("/update").To(UpdatePod))
	podWS.Route(podWS.DELETE("/remove").To(RemovePod))
	podWS.Route(podWS.GET("/getAll").To(GetAllPod))
	container.Add(podWS)

	// service
	serviceWS := new(restful.WebService)
	serviceWS.Path("/services").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	serviceWS.Route(serviceWS.POST("/create").To(CreateService))
	serviceWS.Route(serviceWS.GET("/get").To(GetService))
	serviceWS.Route(serviceWS.POST("/update").To(UpdateService))
	serviceWS.Route(serviceWS.DELETE("/remove").To(RemoveService))
	container.Add(serviceWS)

	// TODO 在此添加新的HTTP请求接口
}
