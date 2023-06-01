package handlers

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/etcd"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
)

func CreateGPUJob(request *restful.Request, response *restful.Response) {
	log.Println("Get create GPUJob request")
	GPUJob := new(object.GPUJob)
	err := request.ReadEntity(&GPUJob)
	if err != nil {
		log.Println(err)
		return
	}
	t := time.Now()
	GPUJob.Metadata.CreationTimestamp = t

	key := "/registry/GPUJobs/default/" + GPUJob.Metadata.Name
	GPUJobString, _ := json.Marshal(GPUJob)
	res := etcd.Put(key, string(GPUJobString))
	response.AddHeader("Content-Type", "text/plain")
	if !res {
		err := response.WriteErrorString(http.StatusNotFound, "GPUJob could not be persisted")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		ret := "ok"
		_, err := response.Write([]byte(ret))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetGPUJob(request *restful.Request, response *restful.Response) {
	GPUJobName := request.PathParameter("name")
	log.Println("Get request: " + GPUJobName)
	key := "/registry/GPUJobs/default/" + GPUJobName
	val := etcd.GetOne(key)
	_, err := response.Write([]byte(val))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func UpdateGPUJob(request *restful.Request, response *restful.Response) {
	log.Println("Get update GPUJob request")
	newGPUJobInfo := object.GPUJob{}
	err := request.ReadEntity(&newGPUJobInfo)
	log.Println(newGPUJobInfo)
	if err != nil {
		log.Println(err)
		return
	}
	newVal, _ := json.Marshal(&newGPUJobInfo)
	key := "/registry/GPUJobs/default/" + newGPUJobInfo.Metadata.Name
	var ret string
	oldValue := etcd.GetOne(key)
	if oldValue == "" {
		ret = "non-existed GPUJob"
		log.Println("update non-existed GPUJob")
		err1 := response.WriteErrorString(500, ret)
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	} else if oldValue == string(newVal) {
		// no update, just return
		ret = "ok"
		_, err1 := response.Write([]byte(ret))
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	} else {
		etcd.Put(key, string(newVal))
		ret = "ok"
		_, err1 := response.Write([]byte(ret))
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	}
}

func RemoveGPUJob(request *restful.Request, response *restful.Response) {
	var rmGPUJobName string
	err := request.ReadEntity(&rmGPUJobName)
	if err != nil {
		return
	}
	log.Println(rmGPUJobName)
	key := "/registry/GPUJobs/default/" + rmGPUJobName
	log.Println("delete key : " + key)
	noError := etcd.Del(key)
	if !noError {
		log.Println("delete GPUJob error")
	}
}

func GetAllGPUJob(request *restful.Request, response *restful.Response) {
	GPUJobMap := etcd.GetDirectory("/registry/GPUJobs")
	msg, _ := json.Marshal(GPUJobMap)
	_, err := response.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
	}
}
