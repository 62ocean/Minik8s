package testFile

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"testing"
)

var serverlessClient = HTTPClient.CreateHTTPClient(global.ServerlessHost)

func TestFunction(t *testing.T) {
	filePath := "test/serverless/hello.py"
	var function object.Function
	function.Name = "test-hello"
	function.Path = filePath
	funjson, _ := json.Marshal(function)
	serverlessClient.Post("/functions/create", funjson)
}

func TestWorkflow(t *testing.T) {
	filePath := "../serverless/testworkflow2.yaml"
	newWf := parseYaml.ParseYaml[object.Workflow](filePath)
	wfJson, _ := json.Marshal(newWf)
	log.Println(newWf)
	serverlessClient.Post("/workflows/create", wfJson)
}
