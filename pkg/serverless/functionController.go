package serverless

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/util/HTTPClient"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
)

type FunctionController interface {
	InitFunction() error

	AddFunction(request *restful.Request, response *restful.Response)
	UpdateFunction(request *restful.Request, response *restful.Response)
	DeleteFunction(request *restful.Request, response *restful.Response)

	TriggerFunction(request *restful.Request, response *restful.Response)
	ExecFunction(funName string, paramsJson string) string
	HoldFunction(function object.RunningFunction)
}

type functionController struct {
	client *HTTPClient.Client

	functionList        map[string]object.Function
	runningFunctionList map[string]object.RunningFunction

	mutex sync.Mutex
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

	// 检查该function是否在运行
	targetFunction, exist := c.runningFunctionList[functionInfo.Name]
	if exist {
		targetFunction.Timer.Stop()
		targetFunction.Timer.Reset(time.Millisecond)
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

	functionName := request.PathParameter("name")

	// 检查该function是否在运行
	targetFunction, exist := c.runningFunctionList[functionName]
	if exist {
		targetFunction.Timer.Stop()
		targetFunction.Timer.Reset(time.Millisecond)
	}

	// 检查该function是否已存在
	_, exist = c.functionList[functionName]
	if !exist {
		fmt.Println("[DELETE FAILED] function " + functionName + " doesn't exist")
		return
	}
	delete(c.functionList, functionName)

	// 需要在docker仓库中删除吗？（会涉及很多细节处理，先不管这个了，运行起来没影响就好）

	c.client.Del("/functions/remove/" + functionName)

	fmt.Println("[DELETE SUCCESSFULLY] function [" + functionName + "] is removed")

}

func (c *functionController) TriggerFunction(request *restful.Request, response *restful.Response) {
	functionName := request.PathParameter("function-name")
	fmt.Print(functionName)
	var params []object.Param
	err := request.ReadEntity(&params)
	if err != nil {
		log.Println("unmarshall paraJson failed")
		return
	}

	paramsjson, _ := json.Marshal(params)
	log.Println("params: " + string(paramsjson))
	ret := c.ExecFunction(functionName, string(paramsjson))

	_, err = response.Write([]byte(ret))
	if err != nil {
		log.Println("write response error")
		return
	}

}

func (c *functionController) ExecFunction(funName string, paramsJson string) string {

	resultJson := "no response"

	c.mutex.Lock()

	targetFunction, exist := c.runningFunctionList[funName]
	if exist {
		// 重置计时器
		log.Println("new request, reset timer to 30s")
		//utils.OutputJson("runningFunction", targetFunction)
		targetFunction.Timer.Stop()
		targetFunction.Timer.Reset(time.Second * 30)
		// 发请求
		resultJson = targetFunction.Client.Post("/", []byte(paramsJson))

		c.mutex.Unlock()

	} else {
		var runningFunction object.RunningFunction
		runningFunction.Function = c.functionList[funName]
		// 起rs
		rs := CreateFunctionRS(funName, c.functionList[funName].Image)
		rsjson, _ := json.Marshal(rs)
		c.client.Post("/replicasets/create", rsjson)
		log.Println("create rs ok")

		time.Sleep(time.Second * 1)

		// 起对应的service并等待
		service := CreateFunctionService(funName)
		runningFunction.ServiceIP = service.Spec.ClusterIP + ":80"
		runningFunction.Client = HTTPClient.CreateHTTPClient("http://" + runningFunction.ServiceIP)
		log.Println("function IP: " + runningFunction.ServiceIP)
		servicejson, _ := json.Marshal(service)
		c.client.Post("/services/create", servicejson)
		log.Println("create service ok")
		for {
			response := c.client.Post("/services/check/function-"+funName, nil)
			var flag int
			_ = json.Unmarshal([]byte(response), &flag)
			if flag == 1 {
				log.Println("service is ready!")
				break
			}
			log.Printf("flag: %d, wait for service ready...\n", flag)
			time.Sleep(time.Second * 1)
		}

		// 起对应的hpa
		hpa := CreateFunctionHPA(funName)
		hpajson, _ := json.Marshal(hpa)
		c.client.Post("/hpas/create", hpajson)
		log.Println("create hpa ok")

		// 开始计时
		runningFunction.Timer = time.NewTimer(time.Second * 30)

		// 添加到内存列表中
		c.runningFunctionList[runningFunction.Function.Name] = runningFunction

		// 发请求
		resultJson = runningFunction.Client.Post("/", []byte(paramsJson))
		c.mutex.Unlock()

		go c.HoldFunction(runningFunction)
	}

	return resultJson
}

func (c *functionController) HoldFunction(function object.RunningFunction) {

	<-function.Timer.C

	// 到时间了
	log.Println("30s passed, no new request, stop running function")
	function.Timer.Stop()

	c.mutex.Lock()

	// 删除service
	serviceName, _ := json.Marshal("function-" + function.Function.Name)
	c.client.Post("/services/remove", serviceName)

	// 删除hpa
	c.client.Del("/hpas/remove/" + "function-" + function.Function.Name)

	// 将目标rs副本设为0（删除所有pod)
	response := c.client.Get("/replicasets/get/" + "function-" + function.Function.Name)
	rs := object.ReplicaSet{}
	_ = json.Unmarshal([]byte(response), &rs)
	rs.Spec.Replicas = 0
	rsjson, _ := json.Marshal(rs)
	c.client.Post("/replicasets/update", rsjson)

	// 删除rs
	c.client.Del("/replicasets/remove/" + "function-" + function.Function.Name)

	delete(c.runningFunctionList, function.Function.Name)

	c.mutex.Unlock()

}

func CreateFunctionRS(functionName string, functionImage string) object.ReplicaSet {
	var rs object.ReplicaSet

	rs.ApiVersion = "apps/v1"
	rs.Kind = "Replicaset"
	rs.Metadata.Name = "function-" + functionName

	rs.Spec.Replicas = 1
	rs.Spec.Selector.MatchLabels.App = "function-" + functionName
	rs.Spec.Selector.MatchLabels.Env = "prod"

	var pod object.Pod
	pod.Metadata.Name = "function-" + functionName
	pod.Metadata.Labels.App = "function-" + functionName
	pod.Metadata.Labels.Env = "prod"

	var container object.Container
	container.Name = "function-" + functionName
	container.Image = functionImage
	container.Ports = append(container.Ports, object.ContainerPort{Port: 8888})
	pod.Spec.Containers = append(pod.Spec.Containers, container)

	rs.Spec.PodTemplate = pod

	return rs
}

func CreateFunctionHPA(functionName string) object.Hpa {
	var hpa object.Hpa

	hpa.ApiVersion = "autoscaling/v2beta2"
	hpa.Kind = "HorizontalPodAutoscaler"
	hpa.Metadata.Name = "function-" + functionName

	hpa.Spec.ScaleTargetRef.Kind = "Replicaset"
	hpa.Spec.ScaleTargetRef.Name = "function-" + functionName

	hpa.Spec.MinReplicas = 1
	hpa.Spec.MaxReplicas = 10

	var metric object.Metric
	metric.Resource.Name = "cpu"
	metric.Resource.Target.AverageUtilization = 0.5
	hpa.Spec.Metrics = append(hpa.Spec.Metrics, metric)

	metric.Resource.Name = "memory"
	metric.Resource.Target.AverageUtilization = 0.5
	hpa.Spec.Metrics = append(hpa.Spec.Metrics, metric)

	return hpa
}

func CreateFunctionService(functionName string) object.Service {
	var newService object.Service
	newService.Kind = "Service"
	newService.Metadata.Name = "function-" + functionName

	rand.Seed(time.Now().Unix())
	newService.Spec.ClusterIP = fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	newService.Spec.Selector.App = "function-" + functionName
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
	c.runningFunctionList = make(map[string]object.RunningFunction)
	err := c.InitFunction()
	if err != nil {
		log.Println("init functions fail")
		return nil
	}

	return c
}
