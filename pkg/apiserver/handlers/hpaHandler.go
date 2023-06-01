package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
	"net/http"
)

func CreateHpa(request *restful.Request, response *restful.Response) {
	log.Println("Get create hpa request")
	hpa := new(object.Hpa)
	err := request.ReadEntity(&hpa)
	if err != nil {
		log.Println(err)
		return
	}
	//uid := hpa.Metadata.Uid
	key := "/registry/hpas/default/" + hpa.Metadata.Name
	rsString, _ := json.Marshal(hpa)
	res := etcd.Put(key, string(rsString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "hpa could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		podQueue := "hpas"
		//err := response.WriteEntity(string(podQueue))
		_, err := response.Write([]byte(podQueue))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetAllHpa(request *restful.Request, response *restful.Response) {
	rsMap := etcd.GetDirectory("/registry/hpas")
	msg, _ := json.Marshal(rsMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetHpa(request *restful.Request, response *restful.Response)    {}
func UpdateHpa(request *restful.Request, response *restful.Response) {}
func RemoveHpa(request *restful.Request, response *restful.Response) {
	hpaName := request.PathParameter("hpaName")
	key := "/registry/hpas/default/" + hpaName
	etcd.Del(key)
}
