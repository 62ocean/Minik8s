package testPod

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"testing"
)

var APIClient = HTTPClient.CreateHTTPClient(global.ServerHost)

func TestGPUJob(t *testing.T) {
	filePath := "../gpuJobAdd.yaml"
	cuFilePath := "../matrixAdd.cu"
	// job存入apiserver
	job := parseYaml.ParseYaml[object.GPUJob](filePath)
	job.Status = object.PENDING
	jobInfo, _ := json.Marshal(job)
	APIClient.Post("/gpuJobs/create", jobInfo)

	// 构造pod 存入apiserver
	port := object.ContainerPort{Port: 8099}
	container := object.Container{
		Name:  "commit_" + "GPUJob_" + job.Metadata.Name,
		Image: "saltfishy/gpu_server:v9",
		Ports: []object.ContainerPort{
			port,
		},
		Command: []string{
			"/apps/main",
		},
		Args: []string{
			job.Metadata.Name,
		},
		CopyFile: cuFilePath,
		CopyDst:  "/apps",
	}
	newPod := object.Pod{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: object.Metadata{
			Name: "GPUJob_" + job.Metadata.Name,
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
	APIClient.Post("/pods/create", podInfo)

	//time.Sleep(time.Second * 30)
	//jobResult := APIClient.Get("/gpuJobs/get/" + job.Metadata.Name)
	//jobR := object.GPUJob{}
	//_ = json.Unmarshal([]byte(jobResult), &jobR)
	//fmt.Println(jobR.Output)
}
