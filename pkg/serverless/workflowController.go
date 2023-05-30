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

type server interface {
	GetFunController() FunctionController
}

type workflowController struct {
	s      server
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
		//initWorkflowMap(&workflow)
		if err != nil {
			fmt.Println("unmarshall workflow failed")
			return err
		}
		c.workflowList[workflow.Metadata.Name] = workflow
	}

	return nil
}

func (c *workflowController) AddWorkflow(request *restful.Request, response *restful.Response) {

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
	initWorkflowMap(&workflowInfo)
	c.workflowList[workflowInfo.Metadata.Name] = workflowInfo
	wfJson, _ := json.Marshal(workflowInfo)

	c.client.Post("/workflows/create", wfJson)

	log.Println("create workflow [" + workflowInfo.Metadata.Name + "] successfully")

}

func initWorkflowMap(workflow *object.Workflow) {
	for _, value := range workflow.Params {
		workflow.ParamsMap[value.Name] = value
	}
	for _, value := range workflow.Steps {
		workflow.StepsMap[value.Name] = value
	}
}

func (c *workflowController) UpdateWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) DeleteWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) GetAllWorkflow(request *restful.Request, response *restful.Response) {

}

func (c *workflowController) TriggerWorkflow(request *restful.Request, response *restful.Response) {
	workflowName := request.PathParameter("workflow-name")
	fmt.Println(workflowName)

	targetWorkflow := c.workflowList[workflowName]

	var paramsJson, retJson string
	pjson, _ := json.Marshal(targetWorkflow.Params)
	paramsJson = string(pjson)
	current := targetWorkflow.Start

	for current != "END" {
		currentStep := targetWorkflow.StepsMap[current]
		if currentStep.Type == "function" {
			// 执行function
			retJson, _ = c.s.GetFunController().ExecFunction(currentStep.Name, paramsJson)
			current = currentStep.Next
		} else if currentStep.Type == "branch" {
			// 解析参数并判断分支
			params := paramsJson2Map(paramsJson)
			retJson = paramsJson
			for _, choice := range currentStep.Choices {
				switch choice.Type {
				case "equal":
					if params[choice.Variable].Value == choice.Value {
						current = choice.Next
						break
					}
				case "notEqual":
					if params[choice.Variable].Value != choice.Value {
						current = choice.Next
						break
					}
				case "moreThan":
					if params[choice.Variable].Value > choice.Value {
						current = choice.Next
						break
					}
				case "lessThan":
					if params[choice.Variable].Value < choice.Value {
						current = choice.Next
						break
					}

				default:

				}
			}
		}
		paramsJson = retJson
	}

	// 最终返回retJson
	_, err := response.Write([]byte(retJson))
	if err != nil {
		log.Println("write to response failed")
		return
	}

}

func paramsJson2Map(paramsJson string) map[string]object.Param {
	var params []object.Param
	err := json.Unmarshal([]byte(paramsJson), &params)
	if err != nil {
		log.Println("unmarshall params failed")
		return nil
	}
	pmap := make(map[string]object.Param)
	for _, value := range params {
		pmap[value.Name] = value
	}
	return pmap
}

func NewWorkflowController(client *HTTPClient.Client, s server) WorkflowController {
	c := &workflowController{}
	c.client = client
	c.workflowList = make(map[string]object.Workflow)
	c.s = s
	err := c.InitWorkflow()
	if err != nil {
		log.Println("init workflows fail")
		return nil
	}

	return c
}
