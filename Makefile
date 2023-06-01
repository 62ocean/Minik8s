GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test

TARGET_KUBELET=kubelet
TARGET_APISERVER=apiserver
TARGET_KUBECTL=kubectl
TARGET_SCHEDULER=scheduler
TARGET_CONTROLLERMANAGER = controllerManager
TARGET_DNS = dns
TARGET_KUBEPROXY = kubeProxy
TARGET_FLANNEL = flannel
.DEFAULT_GOAL := default

GO_TEST_PATH='./test/yaml_test'

# as there is a dir named "test" too, so we need .PHONY to specify this target.
.PHONY:test

all: test master node

build: module apiserver kubectl kubelet scheduler controllerManager dns kubeProxy flannel

master: kubectl apiserver scheduler replicaSet dns dns flannel

node: kubelet kubeProxy flannel

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

kubectl:
	$(GO_BUILD) -o ./build/$(TARGET_KUBECTL) ./cmd/kubectl/main.go

scheduler:
	$(GO_BUILD) -o ./build/$(TARGET_SCHEDULER) ./cmd/scheduler/main.go

kubelet:
	$(GO_BUILD) -o ./build/$(TARGET_KUBELET) ./cmd/kubelet/main.go

controllerManager:
	$(GO_BUILD) -o ./build/$(TARGET_CONTROLLERMANAGER) ./cmd/controllerManager/main.go

dns:
	$(GO_BUILD) -o ./build/$(TARGET_DNS) ./cmd/Dns/main.go

kubeProxy:
	$(GO_BUILD) -o ./build/$(TARGET_KUBEPROXY) ./cmd/kubeProxy/main.go

flannel:
	$(GO_BUILD) -o ./build/$(TARGET_FLANNEL) ./cmd/flannel/main.go

clean:
	rm -rf ./build

master_start:
	sudo /bin/bash -c 'etcd &'
	sudo /bin/bash -c './build/apiserver &'
	sudo /bin/bash -c './build/scheduler &'
	sudo /bin/bash -c './build/controllerManager &'
	sudo /bin/bash -c './build/dns'
	sudo /bin/bash -c './build/coredns &'
	sudo /bin/bash -c './build/flannel &'

node_start:
	sudo /bin/bash -c './build/kubeProxy &'
	sudo /bin/bash -c './build/kubelet -f /builds/520021910279/mini-k8s-2023/utils/templates/node_template.yaml &'
	sudo /bin/bash -c './build/flannel &'

start_all:
	sudo /bin/bash -c 'etcd &'
	sudo /bin/bash -c './build/apiserver &'
	sudo /bin/bash -c './build/scheduler &'
	sudo /bin/bash -c './build/controllerManager &'
	sudo /bin/bash -c './build/kubelet &'
#	sudo /bin/bash -c './build/kubelet -f /builds/520021910279/mini-k8s-2023/utils/templates/node_template.yaml &'
#    sudo /bin/bash -c './build/autoScaler &'
#	sudo /bin/bash -c './build/replicaSet &'
	sudo /bin/bash -c './build/kubeProxy &'
	sudo /bin/bash -c './build/dns'
	sudo /bin/bash -c './build/coredns &'
	sudo /bin/bash -c './build/flannel &'


clean-env:
	sudo /bin/bash -c 'iptables -t nat -F'
	sudo /bin/bash -c 'iptables -t nat -X'
	sudo /bin/bash -c 'systemctl restart docker'
	sudo /bin/bash -c 'etcdctl del "/registry/pods/" --prefix'
	sudo /bin/bash -c 'etcdctl del "/registry/replicasets/" --prefix'
	sudo /bin/bash -c 'etcdctl del "/registry/services/" --prefix'
	sudo /bin/bash -c 'etcdctl del "/registry/endpoints/" --prefix'
	sudo /bin/bash -c 'etcdctl del "/registry/hpas/" --prefix'
	sudo /bin/bash -c 'etcdctl del "/registry/nodes/" --prefix'
	sudo /bin/bash -c 'docker stop $$(docker ps -aq) && docker rm $$(docker ps -aq)'

kill-all:
	-sudo /bin/bash -c 'killall etcd'
	-sudo /bin/bash -c 'killall apiserver'
	-sudo /bin/bash -c 'killall scheduler'
	-sudo /bin/bash -c 'killall controllerManager'
	-sudo /bin/bash -c 'killall kubelet'
	-sudo /bin/bash -c 'killall kubeProxy'
	-sudo /bin/bash -c 'killall dns'
	-sudo /bin/bash -c 'killall coredns'
	-sudo /bin/bash -c 'killall flannel'





