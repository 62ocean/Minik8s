package gpu_server

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/GPU"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"os"
)

func main() {
	// 获取job信息
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	name := os.Args[1]
	response := client.Get("/gpuJobs/get/" + name)
	job := object.GPUJob{}
	_ = json.Unmarshal([]byte(response), &job)
	server := GPU.NewServer(job)
	server.Run()
}
