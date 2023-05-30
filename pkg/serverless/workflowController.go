package serverless

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"k8s/pkg/util/HTTPClient"
	"log"
)

type WorkflowController interface {
	InitWorkflow() error

	AddWorkflow(request *restful.Request, response *restful.Response)
	UpdateWorkflow(request *restful.Request, response *restful.Response)
	DeleteWorkflow(request *restful.Request, response *restful.Response)
	GetAllWorkflow(request *restful.Request, response *restful.Response)

	TriggerWorkflow(request *restful.Request, response *restful.Response)
}

type workflowController struct {
	client *HTTPClient.Client

	workflowList map[string]object.Workflow
}

func (c *workflowController) InitWorkflow() error {
	//得到所有的workflow列表
	response := c.client.Get("/workflows/getAll")
	workflowList := new(map[string]string)
	err := json.Unmarshal([]byte(response), workflowList)
	if err != nil {
		log.Println("unmarshall workflow list failed")
		return err
	}

	// 将所有workflow载入内存
	for _, value := range *workflowList {
		//fmt.Println(value)
		var workflow object.Workflow
		err := json.Unmarshal([]byte(value), &workflow)
		if err != nil {
			fmt.Println("unmarshall workflow failed")
			return err
		}
		c.workflowList[workflow.Metadata.Name] = workflow
	}

	return nil
}

func (c *workflowController) AddWorkflow(request *restful.Request, response *restful.Response) {

	// 拿到函数名字和函数路径
	workflowInfo := object.Workflow{}
	err := request.ReadEntity(&workflowInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}

	// 检查该workflow是否已存在
	_, exist := c.workflowList[workflowInfo.Metadata.Name]
	if exist {
		log.Println("workflow " + workflowInfo.Metadata.Name + " already exist")
		return
	}

	log.Println("start adding workflow " + workflowInfo.Metadata.Name)

	// 向workflowList中添加该workflow, 并将其持久化到etcd中
	c.workflowList[workflowInfo.Metadata.Name] = workflowInfo
	wfJson, _ := json.Marshal(workflowInfo)

	c.client.Post("/workflows/create", wfJson)

	log.Println("create workflow [" + workflowInfo.Metadata.Name + "] successfully")

}

func (c *workflowController) UpdateWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) DeleteWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) GetAllWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) TriggerWorkflow(request *restful.Request, response *restful.Response) {

	// 拿到结果后删除pod（关闭容器）
}

func NewWorkflowController(client *HTTPClient.Client) WorkflowController {
	c := &workflowController{}
	c.client = client
	c.workflowList = make(map[string]object.Workflow)
	err := c.InitWorkflow()
	if err != nil {
		log.Println("init workflows fail")
		return nil
	}

	return c
}
