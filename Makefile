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

VAR ?= default_name

.DEFAULT_GOAL := default

GO_TEST_PATH='./test/yaml_test'

# as there is a dir named "test" too, so we need .PHONY to specify this target.
.PHONY:test

all: test master node

build: module apiserver kubectl kubelet scheduler controllerManager dns kubeProxy flannel

master: kubectl apiserver scheduler replicaSet dns dns flannel

node: kubelet kubeProxy flannel

default: build


testPod: apiserver
	sudo /bin/bash -c './build/apiserver &'

module:
	$(GO_CMD) mod tidy

apiserver:
	$(GO_BUILD) -o ./build/$(TARGET_APISERVER) ./cmd/apiserver/main.go

kubectl:
	$(GO_BUILD) -o ./$(TARGET_KUBECTL) ./cmd/kubectl/main.go

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
	sudo /bin/bash -c 'etcd -listen-client-urls="http://192.168.1.6:2379,http://localhost:2379" -advertise-client-urls="http://192.168.1.6:2379"   &'
	sudo /bin/bash -c './build/apiserver &'
	sleep 5
	sudo /bin/bash -c './build/scheduler &'
	sudo /bin/bash -c './build/controllerManager &'
	sudo /bin/bash -c './build/dns'
	sudo /bin/bash -c './build/coredns &'
	sudo /bin/bash -c './build/flannel &'

node_start:
	echo "$(VAR)"
	sudo /bin/bash -c './build/kubeProxy &'
	sudo /bin/bash -c './build/kubelet $(VAR) &'
	sudo /bin/bash -c './build/flannel &'

start_all:
	sudo /bin/bash -c 'etcd -listen-client-urls="http://192.168.1.6:2379,http://localhost:2379" -advertise-client-urls="http://192.168.1.6:2379"   &'
	sudo /bin/bash -c './build/apiserver &'
	sleep 5
	sudo /bin/bash -c './build/scheduler &'
	sudo /bin/bash -c './build/controllerManager &'
	sudo /bin/bash -c './build/kubelet $(VAR) &'
	sudo /bin/bash -c './build/kubeProxy &'
	sudo /bin/bash -c './build/dns'
	sudo /bin/bash -c './build/coredns &'
	sudo /bin/bash -c './build/flannel &'


clean-env:
	-sudo /bin/bash -c 'iptables -t nat -F'
	-sudo /bin/bash -c 'iptables -t nat -X'
	-sudo /bin/bash -c 'systemctl restart docker'
	-sudo /bin/bash -c 'docker stop $$(docker ps -aq) && docker rm $$(docker ps -aq)'
	-sudo /bin/bash -c 'etcdctl del "" --prefix'

kill-all:
	-sudo /bin/bash -c 'ps -ef | grep etcd | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall etcd'
	-sudo /bin/bash -c 'ps -ef | grep apiserver | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall apiserver'
	-sudo /bin/bash -c 'ps -ef | grep scheduler | awk '{print $2}' | xargs kill -9 '
	-sudo /bin/bash -c 'killall scheduler'
	-sudo /bin/bash -c 'ps -ef | grep controllerManager | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall controllerManager'
	-sudo /bin/bash -c 'ps -ef | grep kubelet | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall kubelet'
	-sudo /bin/bash -c 'ps -ef | grep kubeProxy | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall kubeProxy'
	-sudo /bin/bash -c 'ps -ef | grep dns | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall dns'
	-sudo /bin/bash -c 'ps -ef | grep coredns | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall coredns'
	-sudo /bin/bash -c 'ps -ef | grep flannel | awk '{print $2}' | xargs kill -9'
	-sudo /bin/bash -c 'killall flannel'





