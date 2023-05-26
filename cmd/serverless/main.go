package main

import (
	"k8s/pkg/serverless"
)

func main() {
	server, _ := serverless.CreateAPIServer()
	server.StartServer()
}
