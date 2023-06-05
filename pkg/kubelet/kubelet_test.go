package kubelet

import (
	"k8s/object"
	"k8s/pkg/kubelet/cache"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"sync"
	"testing"
)

func TestKubelet_createPod(t *testing.T) {
	type fields struct {
		client        *HTTPClient.Client
		node          object.Node
		podSubscriber *subscriber.Subscriber
		podQueue      string
		podHandler    podHandler
		pods          map[string]*cache.PodCache
		toBeDel       map[string]*cache.PodCache
		mutex         sync.Mutex
	}
	type args struct {
		podInfo *cache.PodCache
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kub := &Kubelet{
				client:        tt.fields.client,
				node:          tt.fields.node,
				podSubscriber: tt.fields.podSubscriber,
				podQueue:      tt.fields.podQueue,
				podHandler:    tt.fields.podHandler,
				pods:          tt.fields.pods,
				toBeDel:       tt.fields.toBeDel,
				mutex:         tt.fields.mutex,
			}
			kub.createPod(tt.args.podInfo)
		})
	}
}
