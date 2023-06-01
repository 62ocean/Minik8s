GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test

TARGET_KUBELET=kubelet
TARGET_APISERVER=apiserver
TARGET_KUBECTL=kubectl
TARGET_SCHEDULER=scheduler
TARGET_CONTROLLERMANAGER = controllerManager
.DEFAULT_GOAL := default

GO_TEST_PATH='./test/yaml_test'

# as there is a dir named "test" too, so we need .PHONY to specify this target.
.PHONY:test

all: test master node

build: module apiserver kubectl kubelet scheduler controllerManager

master: kubectl apiserver scheduler replicaSet

node: kubelet

default: build


test1:
	go test -v ./test/yaml_test/yaml_test.go
	go test -v ./test/etcd_test/etcd_test.go
	go test -v ./test/container_test/container_test.go
	go test -v ./test/node_test/node_test.go
test:
	go test -v ./test/pod_test/pod_test.go
	go test -v ./test/node_test/node1_test.go
	go test -v ./test/service_test/service_test.go
	go test -v ./test/auto_test/auto_test.go
	# go test -v ./test/replicaSet_test/replicaSet_test.go

module:
	$(GO_CMD) mod tidy

apiserver:
	$(GO_BUILD) -o ./build/$(TARGET_APISERVER) ./cmd/apiserver/main.go

kubectl: apiserver
	$(GO_BUILD) -o ./build/$(TARGET_KUBECTL) ./cmd/kubectl/main.go

scheduler: apiserver
	$(GO_BUILD) -o ./build/$(TARGET_SCHEDULER) ./cmd/scheduler/main.go

kubelet: apiserver
	$(GO_BUILD) -o ./build/$(TARGET_KUBELET) ./cmd/kubelet/main.go

controllerManager: apiserver
	$(GO_BUILD) -o ./build/$(TARGET_CONTROLLERMANAGER) ./cmd/controllerManager/main.go

clean:
	rm -rf ./build

master_start:
#	sudo ./build/apiserver &
#	sudo ./build/scheduler &
#	sudo ./build/kubectl
	sudo /bin/bash -c './build/apiserver &'
	sudo /bin/bash -c 'sleep 5'
	sudo /bin/bash -c './build/scheduler &'
	sudo /bin/bash -c './build/autoScaler &'
	sudo /bin/bash -c './build/replicaSet &'
#	sudo sh -c './build/kubectl &'

node_start:
#	sudo ./build/kubeproxy &
#	sudo ./build/kubelet
	sudo /bin/bash -c './build/kubeproxy &'
#	sudo ./build/kubelet -f ./utils/templates/node_template.yaml
	sudo /bin/bash -c './build/kubelet -f /builds/520021910279/mini-k8s-2023/utils/templates/node_template.yaml &'
