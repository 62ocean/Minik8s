package kubelet

import "k8s/pkg/kubelet/HTTPClient"

// TODO:记得修改host到对应ip
const serverHost = "localhost:8080"

type Kubelet struct {
	client *HTTPClient.Client
}

// Run kubelet运行的入口函数
func (kl *Kubelet) Run() {

}

// NewKubelet kubelet对象的构造函数
func NewKubelet() (*Kubelet, error) {
	client := HTTPClient.CreateHTTPClient(serverHost)
	// TODO 构建node对象传递到APIserver处
	client.Post("/node/create")
	return nil, nil
}
