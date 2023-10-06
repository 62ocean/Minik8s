package testPod

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/util/parseYaml"
	"log"
	"testing"
)

func TestPod(t *testing.T) {
	filePath := "../podConfigTest.yaml"
	log.Println("create pod: ", filePath)
	newPod := parseYaml.ParseYaml[object.Pod](filePath)
	podJson, _ := json.Marshal(newPod)
	log.Println(newPod)
	APIClient.Post("/pods/create", podJson)
}
