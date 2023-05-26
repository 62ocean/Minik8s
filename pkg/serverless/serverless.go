package serverless

import (
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)

/*-----------------APIServer-----------------*/

type APIServer struct {
	wsContainer *restful.Container
	c           FunctionController
}

// CreateAPIServer 初始化APIServer结构体中的内容
func CreateAPIServer() (*APIServer, error) {

	// HTTP server
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	// construct APIServer
	server := APIServer{
		wsContainer: wsContainer,
		c:           &functionController{},
	}

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
	invokeWS.Route(invokeWS.POST("/{function-name}").To(s.c.TriggerFunction))
	s.wsContainer.Add(invokeWS)

	functionWS := new(restful.WebService)
	functionWS.Path("/functions").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	functionWS.Route(functionWS.POST("/create").To(s.c.AddFunction))
	functionWS.Route(functionWS.POST("/update").To(s.c.UpdateFunction))
	functionWS.Route(functionWS.POST("/remove").To(s.c.DeleteFunction))
	functionWS.Route(functionWS.POST("/getAll").To(s.c.GetAllFunction))
	s.wsContainer.Add(functionWS)
}
