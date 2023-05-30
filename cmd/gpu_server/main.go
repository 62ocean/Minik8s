package main

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/GPU"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"os"
)

func main() {
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	name := os.Args[1]
	response := client.Get("/gpuJobs/get/" + name)
	job := object.GPUJob{}
	_ = json.Unmarshal([]byte(response), &job)
	server := GPU.NewServer(job)
	server.Run()
}
