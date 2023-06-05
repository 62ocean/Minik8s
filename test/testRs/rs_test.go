package testFile

import (
	"encoding/json"
	"k8s/object"
	"k8s/pkg/util/parseYaml"
	"log"
	"testing"
)

func TestRS(t *testing.T) {
	filePath := "../ReplicasetConfigTest.yaml"
	log.Println("create rs: ", filePath)
	newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
	rsJson, _ := json.Marshal(newRS)
	log.Println(newRS)
	APIClient.Post("/replicasets/create", rsJson)
}
