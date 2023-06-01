package serverless

import (
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"log"
	"net/http"
)

/*-----------------APIServer-----------------*/

type APIServer struct {
	wsContainer   *restful.Container
	funController FunctionController
	wfController  WorkflowController

	client *HTTPClient.Client
}

// CreateAPIServer 初始化APIServer结构体中的内容
func CreateAPIServer() (*APIServer, error) {

	// HTTP server
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	// construct APIServer
	server := APIServer{
		wsContainer: wsContainer,
	}
	server.client = HTTPClient.CreateHTTPClient(global.ServerHost)
	server.funController = NewFunctionController(server.client)
	server.wfController = NewWorkflowController(server.client, &server)

	return &server, nil
}

// StartServer 开始监听相关端口请求
func (s *APIServer) StartServer() {

	s.InitWebServer()

	server := &http.Server{Addr: ":8090", Handler: s.wsContainer}
	defer func(server *http.Server) {
		err := server.Close()
		if err != nil {
			return
		}
	}(server)
	log.Fatal(server.ListenAndServe())
}

func (s *APIServer) InitWebServer() {
	invokeWS := new(restful.WebService)
	invokeWS.Path("/invoke").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	invokeWS.Route(invokeWS.POST("/function/{function-name}").To(s.funController.TriggerFunction))
	invokeWS.Route(invokeWS.POST("/workflow/{workflow-name}").To(s.wfController.TriggerWorkflow))
	s.wsContainer.Add(invokeWS)

	functionWS := new(restful.WebService)
	functionWS.Path("/functions").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	functionWS.Route(functionWS.POST("/create").To(s.funController.AddFunction))
	functionWS.Route(functionWS.POST("/update").To(s.funController.UpdateFunction))
	functionWS.Route(functionWS.DELETE("/remove/{name}").To(s.funController.DeleteFunction))
	functionWS.Route(functionWS.GET("/getAll").To(s.funController.GetAllFunction))
	s.wsContainer.Add(functionWS)

	workflowWS := new(restful.WebService)
	workflowWS.Path("/workflows").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	workflowWS.Route(workflowWS.POST("/create").To(s.wfController.AddWorkflow))
	workflowWS.Route(workflowWS.POST("/update").To(s.wfController.UpdateWorkflow))
	workflowWS.Route(workflowWS.POST("/remove").To(s.wfController.DeleteWorkflow))
	workflowWS.Route(workflowWS.POST("/getAll").To(s.wfController.GetAllWorkflow))
	s.wsContainer.Add(workflowWS)
}

func (s *APIServer) GetFunController() FunctionController {
	return s.funController
}
