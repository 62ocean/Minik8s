package serverless

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"k8s/object"
	"k8s/pkg/util/HTTPClient"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type FunctionController interface {
	InitFunction() error

	AddFunction(request *restful.Request, response *restful.Response)
	UpdateFunction(request *restful.Request, response *restful.Response)
	DeleteFunction(request *restful.Request, response *restful.Response)
	GetAllFunction(request *restful.Request, response *restful.Response)

	TriggerFunction(request *restful.Request, response *restful.Response)
	ExecFunction(funName string, paramsJson string) (string, error)
}

type functionController struct {
	client *HTTPClient.Client

	functionList map[string]object.Function
}

func (c *functionController) InitFunction() error {
	//得到所有的function列表
	response := c.client.Get("/functions/getAll")
	functionList := new(map[string]string)
	err := json.Unmarshal([]byte(response), functionList)
	if err != nil {
		log.Println("unmarshall function list failed")
		return err
	}

	// 将所有function载入内存
	for _, value := range *functionList {
		//fmt.Println(value)
		var function object.Function
		err := json.Unmarshal([]byte(value), &function)
		if err != nil {
			fmt.Println("unmarshall function failed")
			return err
		}
		c.functionList[function.Name] = function
	}

	return nil
}

func buildAndPushImage(functionInfo object.Function) error {
	// 生成对应的Dockerfile
	filedir := filepath.Dir(functionInfo.Path)
	filename := filepath.Base(functionInfo.Path)

	dockerfilePath := filedir + "/Dockerfile"
	fmt.Println(dockerfilePath)

	dockerfileData := "FROM python:3.11\n"
	dockerfileData += "WORKDIR ./" + functionInfo.Name + "\n"
	dockerfileData += "ADD . .\n"
	dockerfileData += "RUN pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple\n"
	dockerfileData += "EXPOSE 8888\n"
	dockerfileData += "CMD [\"python\", \"./" + filename + "\"]\n"

	file, err := os.Create(dockerfilePath)
	defer file.Close()
	if err != nil {
		log.Println("create dockerfile failed")
		return err
	}
	_, _ = file.WriteString(dockerfileData)

	// 创建容器镜像并将其推送至dockerhub
	cmd := exec.Command("bash", "pkg/serverless/buildImage.sh",
		filedir, functionInfo.Image, functionInfo.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return nil
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

	// 检查该function是否已存在
	_, exist := c.functionList[functionInfo.Name]
	if exist {
		fmt.Println("[CREATE FAILED] function " + functionInfo.Name + " already exist")
		return
	}

	log.Println("start adding function " + functionInfo.Name + " from " + functionInfo.Path)

	// build and push
	functionInfo.Version = 0
	functionInfo.ImageName = strings.ToLower("ocean62/" + functionInfo.Name + "-" + uuid.New().String())
	functionInfo.Image = functionInfo.ImageName + ":v" + strconv.Itoa(functionInfo.Version)
	err = buildAndPushImage(functionInfo)
	if err != nil {
		log.Println("build and push image failed")
		return
	}

	// 向functionList中添加该function, 并将其持久化到etcd中
	c.functionList[functionInfo.Name] = functionInfo
	funJson, _ := json.Marshal(functionInfo)

	c.client.Post("/functions/create", funJson)

	fmt.Println("[CREATE SUCCESSFULLY] function [" + functionInfo.Name + "] is available now")

}

func (c *functionController) UpdateFunction(request *restful.Request, response *restful.Response) {
	functionInfo := object.Function{}
	err := request.ReadEntity(&functionInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}

	// 检查该function是否已存在
	function, exist := c.functionList[functionInfo.Name]
	if !exist {
		fmt.Println("[UPDATE FAILED] function " + functionInfo.Name + " doesn't exist")
		return
	}

	function.Path = functionInfo.Path
	function.Version++
	function.Image = function.ImageName + ":v" + strconv.Itoa(function.Version)

	// build and push
	err = buildAndPushImage(function)
	if err != nil {
		log.Println("build and push image failed")
		return
	}

	// 向functionList中更改该function, 并将其持久化到etcd中
	c.functionList[function.Name] = function
	funJson, _ := json.Marshal(function)

	c.client.Post("/functions/update", funJson)

	fmt.Println("[UPDATE SUCCESSFULLY] function [" + functionInfo.Name + "] is updated")

}

func (c *functionController) DeleteFunction(request *restful.Request, response *restful.Response) {
	functionInfo := object.Function{}
	err := request.ReadEntity(&functionInfo)
	//fmt.Println(newRSInfo)
	if err != nil {
		log.Println(err)
		return
	}

	// 检查该function是否已存在
	_, exist := c.functionList[functionInfo.Name]
	if !exist {
		fmt.Println("[DELETE FAILED] function " + functionInfo.Name + " doesn't exist")
		return
	}

	// 需要在docker仓库中删除吗？（会涉及很多细节处理，先不管这个了，运行起来没影响就好）

	delete(c.functionList, functionInfo.Name)

	fmt.Println("[DELETE SUCCESSFULLY] function [" + functionInfo.Name + "] is removed")

}

func (c *functionController) GetAllFunction(request *restful.Request, response *restful.Response) {

}

func (c *functionController) TriggerFunction(request *restful.Request, response *restful.Response) {
	functionName := request.PathParameter("function-name")
	fmt.Print(functionName)
	var paramsJson string
	_ = request.ReadEntity(&paramsJson)

	// 检查该pod是否存在，如不存在，创建pod

	// 向etcd中添加一个pod
	newPod := CreateFunctionPod(functionName, c.functionList[functionName].Image)
	podJson, _ := json.Marshal(newPod)
	c.client.Post("/pods/create", podJson)

	//向etcd中添加一个service
	newService := CreateFunctionService(functionName)
	serviceJson, _ := json.Marshal(newService)
	c.client.Post("/services/create", serviceJson)

	//执行函数
	ret, _ := c.ExecFunction(functionName, paramsJson)
	_, err := response.Write([]byte(ret))
	if err != nil {
		log.Println("write to response failed")
		return
	}

	// 拿到结果后删除pod（关闭容器）
}

func (c *functionController) ExecFunction(funName string, paramsJson string) (string, error) {

	//var params []object.Param
	//err := json.Unmarshal([]byte(paramsJson), &params)
	//if err != nil {
	//	log.Println("unmarshall params failed")
	//	return "", err
	//}
	//workflowInfo := object.Workflow{}
	//err :=
	////fmt.Println(newRSInfo)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	// 把params作为post内容发给对应的容器端口，得到返回结果string
	return "111", nil
}

func CreateFunctionPod(functionName string, functionImage string) object.Pod {
	var pod object.Pod

	pod.ApiVersion = "v1"
	pod.Kind = "Pod"

	pod.Metadata.Uid = uuid.New().String()
	pod.Metadata.Name = "function-" + functionName + "-" + pod.Metadata.Uid
	pod.Metadata.Labels.App = functionName
	pod.Metadata.Labels.Env = "prod"

	var container object.Container
	container.Name = "function-" + functionName + "-" + pod.Metadata.Uid
	container.Image = functionImage
	container.Ports = append(container.Ports, object.ContainerPort{Port: 8888})

	pod.Spec.Containers = append(pod.Spec.Containers, container)

	return pod
}

func CreateFunctionService(functionName string) object.Service {
	var newService object.Service
	newService.Kind = "Service"
	newService.Metadata.Name = "service-" + functionName
	newService.Spec.ClusterIP = "10.10.10.10"
	newService.Spec.Selector.App = functionName
	newService.Spec.Selector.Env = "prod"
	var port object.ServicePort
	port.Protocol = "TCP"
	port.Port = 80
	port.TargetPort = 8888
	newService.Spec.Ports = append(newService.Spec.Ports, port)

	return newService
}

func NewFunctionController(client *HTTPClient.Client) FunctionController {
	c := &functionController{}
	c.client = client
	c.functionList = make(map[string]object.Function)
	err := c.InitFunction()
	if err != nil {
		log.Println("init functions fail")
		return nil
	}

	return c
}
