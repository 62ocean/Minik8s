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
	nodeWS.Route(nodeWS.GET("/getAll").To(GetAllNode))
	container.Add(nodeWS)

	// pod
	podWS := new(restful.WebService)
	podWS.Path("/pods").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	podWS.Route(podWS.POST("/create").To(CreatePod))
	podWS.Route(podWS.GET("/get").To(GetPod))
	podWS.Route(podWS.POST("/update").To(UpdatePod))
	podWS.Route(podWS.POST("/remove").To(RemovePod))
	podWS.Route(podWS.GET("/getAll").To(GetAllPod))
	container.Add(podWS)

	// replicaset
	replicasetWS := new(restful.WebService)
	replicasetWS.Path("/replicasets").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	replicasetWS.Route(replicasetWS.POST("/create").To(CreateReplicaset))
	replicasetWS.Route(replicasetWS.GET("/get").To(GetReplicaset))
	replicasetWS.Route(replicasetWS.POST("/update").To(UpdateReplicaset))
	replicasetWS.Route(replicasetWS.DELETE("/remove/{ip}").To(RemoveReplicaset))
	replicasetWS.Route(replicasetWS.GET("/getAll").To(GetAllReplicaset))
	container.Add(replicasetWS)

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

	// hpa
	hpaWS := new(restful.WebService)
	hpaWS.Path("/hpas").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	hpaWS.Route(hpaWS.POST("/create").To(CreateHpa))
	hpaWS.Route(hpaWS.GET("/get").To(GetHpa))
	hpaWS.Route(hpaWS.POST("/update").To(UpdateHpa))
	hpaWS.Route(hpaWS.DELETE("/remove").To(RemoveHpa))
	hpaWS.Route(hpaWS.GET("/getAll").To(GetAllHpa))
	container.Add(hpaWS)

	// TODO 在此添加新的HTTP请求接口
}
