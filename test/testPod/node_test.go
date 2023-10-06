package testPod

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	response := APIClient.Get("/nodes/getAll")
	var nodeList map[string]string
	_ = json.Unmarshal([]byte(response), &nodeList)
	fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
	for _, val := range nodeList {
		nodeStorage := object.NodeStorage{}
		_ = json.Unmarshal([]byte(val), &nodeStorage)
		createTime := nodeStorage.Node.Metadata.CreationTimestamp
		newtime := time.Now()
		d := newtime.Sub(createTime)
		fmt.Printf("%s\t\t\t%s\t\t\t%s\n", nodeStorage.Node.Metadata.Name, nodeStorage.Status.ToString(), d.Truncate(time.Second).String())
	}
}
