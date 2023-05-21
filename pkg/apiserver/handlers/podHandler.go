package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/pkg/etcd"
)

func CreatePod(request *restful.Request, response *restful.Response) {
}

func GetPod(request *restful.Request, response *restful.Response) {}

//	func UpdatePod(request *restful.Request, response *restful.Response) {
//		newPodInfo := object.PodStorage{}
//		err := request.ReadEntity(newPodInfo)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		newVal, _ := json.Marshal(newPodInfo)
//		key := "/registry/pods/" + newPodInfo
//		etcd.Put()
//
// }
func RemovePod(request *restful.Request, response *restful.Response) {}
func GetAllPod(request *restful.Request, response *restful.Response) {
	podMap := etcd.GetDirectory("/registry/pods")
	msg, _ := json.Marshal(podMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}
