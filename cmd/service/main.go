package main

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
)

func main() {

	service := parseYaml.ParseYaml[object.Service]("./serviceConfigTest.yaml")

	bytes, _ := json.Marshal(service)

	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	client.Post("/services/create", bytes)

	fmt.Println("here")
}
