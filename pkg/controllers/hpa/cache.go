package hpa

import (
	"encoding/json"
	"fmt"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"log"
	"time"
)

type Cache interface {
	UpdatePodStatus()
	SyncLoop()

	GetPodStatusList() map[string]string
}

type cache struct {
	podList map[string]string

	client *HTTPClient.Client
}

func (cache *cache) UpdatePodStatus() {
	response := cache.client.Get("/pods/getAll")
	var podList map[string]string
	err := json.Unmarshal([]byte(response), &podList)
	if err != nil {
		fmt.Println("[hpa cache] unmarshall podlist failed")
		return
	}
	cache.podList = podList
}

func (cache *cache) SyncLoop() {

	cache.client = HTTPClient.CreateHTTPClient(global.ServerHost)

	//每隔10s同步一次pod状态
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		log.Println("[hpa cache] update pod status")
		cache.UpdatePodStatus()
		fmt.Println(cache.podList)
	}
}

func (cache *cache) GetPodStatusList() map[string]string {
	return cache.podList
}

func NewCache() Cache {
	c := &cache{}

	return c
}
