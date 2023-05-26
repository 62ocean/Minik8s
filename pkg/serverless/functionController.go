package serverless

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type FunctionController interface {
	AddFunction(request *restful.Request, response *restful.Response)
	UpdateFunction(request *restful.Request, response *restful.Response)
	DeleteFunction(request *restful.Request, response *restful.Response)
	GetAllFunction(request *restful.Request, response *restful.Response)

	TriggerFunction(request *restful.Request, response *restful.Response)
}

type functionController struct {
}

func (c *functionController) AddFunction(request *restful.Request, response *restful.Response) {
	// 拿到函数名字和函数路径
	functionInfo := object.Function{}
	err := request.ReadEntity(&functionInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}

	// 生成对应的Dockerfile
	filedir := filepath.Dir(functionInfo.Path)
	filename := filepath.Base(functionInfo.Path)

	dockerfilePath := filedir + "/Dockerfile"
	fmt.Println(dockerfilePath)

	dockerfileData := "FROM python:3.11\n"
	dockerfileData += "WORKDIR ./" + functionInfo.Name + "\n"
	dockerfileData += "ADD . .\n"
	dockerfileData += "RUN pip install -r requirements.txt\n"
	dockerfileData += "EXPOSE 8888\n"
	dockerfileData += "CMD [\"python\", \"./" + filename + "\"]\n"

	file, err := os.Create(dockerfilePath)
	defer file.Close()
	if err != nil {
		fmt.Println("create dockerfile failed")
		return
	}
	_, _ = file.WriteString(dockerfileData)

	// 生成对应的requirements.txt
	os.Chdir(filedir) //最后要换回来

	cmd := exec.Command("ls")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()

	cmd = exec.Command("pipreqs", ".", "--encodin", "utf8", "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// build

	// push
}

// ---------- tool functions -----------------

func (c *functionController) UpdateFunction(request *restful.Request, response *restful.Response) {

}

func (c *functionController) DeleteFunction(request *restful.Request, response *restful.Response) {

}

func (c *functionController) GetAllFunction(request *restful.Request, response *restful.Response) {

}

func (c *functionController) TriggerFunction(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("function-name")
	fmt.Print(name)
}

//type Controller interface {
//	Start()
//	FunctionInit() error
//	//FunctionChangeHandler(eventType object.EventType, rs object.Function)
//	AddFunction(fun object.Function)
//	DeleteFunction(fun object.Function)
//	UpdateFunction(fun object.Function)
//}
//
//type controller struct {
//	client  *HTTPClient.Client
//	s       *subscriber.Subscriber
//	handler *functionHandler
//}
//
//func (c *controller) Start() {
//	err := c.s.Subscribe("functions", subscriber.Handler(c.handler))
//	if err != nil {
//		fmt.Println("[function controller] subscribe function failed")
//		return
//	}
//}
//
//func (c *controller) FunctionInit() error {
//	return nil
//}
//
//func (c *controller) AddFunction(fun object.Function) {
//	log.Print("[function controllers] create function: " + fun.Name)
//}
//
//func (c *controller) DeleteFunction(fun object.Function) {
//	log.Print("[function controllers] delete function: " + fun.Name)
//}
//
//func (c *controller) UpdateFunction(fun object.Function) {
//
//	log.Print("[function controllers] update function: " + fun.Name)
//
//	//fmt.Println(fun.Code)
//	// 假设只会更改function内容，不会更改名字
//}
//
//func NewController() Controller {
//	c := &controller{}
//	c.client = HTTPClient.CreateHTTPClient(global.ServerHost)
//
//	//初始化当前etcd中的function
//	err := c.FunctionInit()
//	if err != nil {
//		fmt.Println("[function controller] function init failed")
//		return nil
//	}
//
//	//创建subscribe监听function的变化
//	c.s, _ = subscriber.NewSubscriber(global.MQHost)
//	c.handler = &functionHandler{
//		c: c,
//	}
//
//	return c
//}
//
//// --------------------function change handler----------------
//
//type functionHandler struct {
//	c *controller
//}
//
//func (h *functionHandler) Handle(msg []byte) {
//	log.Println("[function controller] receive function change msg")
//
//	var msgObject object.MQMessage
//	err := json.Unmarshal(msg, &msgObject)
//	if err != nil {
//		fmt.Println("[function controller] unmarshall msg failed")
//		return
//	}
//
//	var function object.Function
//	if msgObject.EventType == object.DELETE {
//		err = json.Unmarshal([]byte(msgObject.PrevValue), &function)
//	} else {
//		err = json.Unmarshal([]byte(msgObject.Value), &function)
//	}
//
//	if err != nil {
//		fmt.Println("[function controller] unmarshall changed function failed")
//		return
//	}
//
//	switch msgObject.EventType {
//	case object.CREATE:
//		h.c.AddFunction(function)
//	case object.DELETE:
//		h.c.DeleteFunction(function)
//	case object.UPDATE:
//		h.c.UpdateFunction(function)
//	}
//}
