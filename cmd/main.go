package main

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"os"
	"runtime"
)

func init() {
	sysType := runtime.GOOS
	var logFile *os.File
	var err error
	if sysType == "linux" || sysType == "darwin" {
		// LINUX系统或者MAC
		logFile, err = os.OpenFile("../log/main.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if sysType == "windows" {
		// windows系统
		logFile, err = os.OpenFile("log/main.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	}
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	//etcd.EtcdTest()
	//apiserver.StartServer()
	//kubectl.CmdExec()
	//fmt.Println("hello world")
	//log.Println("test Log!")

	/*--------------------------KUBECTL FOR GPU-----------------------------*/
	// job存入apiserver
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	job := parseYaml.ParseYaml[object.GPUJob]("../../test/gpuJobAdd.yaml")
	job.Status = object.PENDING
	jobInfo, _ := json.Marshal(job)
	client.Post("/gpuJobs/create", jobInfo)

	// 构造pod 存入apiserver
	port := object.ContainerPort{Port: 8080}
	container := object.Container{
		Name:  "commit_" + "CPUJob_" + job.Metadata.Name,
		Image: "saltfishy/gpu_server:latest",
		Ports: []object.ContainerPort{
			port,
		},
		Command: []string{
			"/bin/sh",
		},
		Args: []string{
			"-c",
			"./gpu_server " + job.Metadata.Name,
		},
	}
	newPod := object.Pod{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: object.Metadata{
			Name: "CPUJob_" + job.Metadata.Name,
			Labels: object.Labels{
				App: "GPU",
				Env: "prod",
			},
		},
		Spec: object.PodSpec{
			Containers: []object.Container{
				container,
			},
		},
	}
	podInfo, _ := json.Marshal(newPod)
	client.Post("/pods/create", podInfo)
	/*--------------------------KUBECTL FOR GPU-----------------------------*/
}
