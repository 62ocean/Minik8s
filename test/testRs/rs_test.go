package testRs

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"testing"
)

var APIClient = HTTPClient.CreateHTTPClient(global.ServerHost)

func TestRS(t *testing.T) {
	filePath := "../ReplicasetConfigTest.yaml"
	log.Println("create rs: ", filePath)
	newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
	rsJson, _ := json.Marshal(newRS)
	log.Println(newRS)
	APIClient.Post("/replicasets/create", rsJson)
}
