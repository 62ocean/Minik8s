package testHPA

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

func TestHPA(t *testing.T) {
	filePath := "../hpaConfigTest.yaml"
	log.Println("create hpa: ", filePath)
	newHpa := parseYaml.ParseYaml[object.Hpa](filePath)
	hpaJson, _ := json.Marshal(newHpa)
	log.Println(newHpa)
	APIClient.Post("/hpas/create", hpaJson)
}
