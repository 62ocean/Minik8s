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
	serviceWS.Route(serviceWS.POST("/get").To(GetService))
	serviceWS.Route(serviceWS.GET("/getAll").To(GetAllService))
	serviceWS.Route(serviceWS.POST("/update").To(UpdateService))
	serviceWS.Route(serviceWS.DELETE("/remove").To(RemoveService))
	serviceWS.Route(serviceWS.POST("/check/{serviceName}").To(CheckService))
	container.Add(serviceWS)

	// endpoint
	endpointWS := new(restful.WebService)
	endpointWS.Path("/endpoints").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	endpointWS.Route(endpointWS.POST("/get").To(GetEndpoint))
	container.Add(endpointWS)

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

	// Dns
	dnsWS := new(restful.WebService)
	dnsWS.Path("/dns").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	dnsWS.Route(dnsWS.POST("/create").To(CreateDns))
	dnsWS.Route(dnsWS.GET("/get").To(GetDns))

	container.Add(dnsWS)

	//GPUJob
	GPUJobWS := new(restful.WebService)
	GPUJobWS.Path("/gpuJobs").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	GPUJobWS.Route(GPUJobWS.POST("/create").To(CreateGPUJob))
	GPUJobWS.Route(GPUJobWS.GET("/get/{name}").To(GetGPUJob))
	GPUJobWS.Route(GPUJobWS.POST("/update").To(UpdateGPUJob))
	GPUJobWS.Route(GPUJobWS.POST("/remove").To(RemoveGPUJob))
	GPUJobWS.Route(GPUJobWS.GET("/getAll").To(GetAllGPUJob))
	container.Add(GPUJobWS)

	// TODO 在此添加新的HTTP请求接口
}
