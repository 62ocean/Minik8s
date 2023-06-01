package main

import "k8s/pkg/kubeProxy"

func main() {
	proxy := kubeProxy.CreateKubeProxy()
	proxy.Run()
}
