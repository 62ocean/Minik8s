package kubeProxy

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/msgQueue/subscriber"
	"log"
	"os/exec"
)

func RunCommand(cmd string) string {
	fmt.Printf("RunCmd: %s\n", cmd)
	command := exec.Command("/bin/bash", "-c", cmd)
	// out, err := command.Output()
	// if err != nil {
	// 	fmt.Println("output : ")
	// 	fmt.Println(out)
	// }
	out, err := command.CombinedOutput()
	fmt.Printf("out: %s", string(out))
	if err != nil {
		// fmt.Printf("ERROR: run cmd error: %s\n", err.Error())
		// if err.Error() != "exit status 1" && err.Error() != "exit status 4" {
		// 	panic("ERROR: " + err.Error())
		// }
	}
	return string(out)
}

func ensureRuleAppend(cmd string) {
	out := RunCommand("iptables -w  -t nat -C" + " " + cmd)
	if out != "" {
		RunCommand("iptables -w -t nat -A" + " " + cmd)

	}
}

func ensureRuleDelete(cmd string) {
	out := RunCommand("iptables -w -t nat -C" + " " + cmd)
	if out == "" {
		RunCommand("iptables -w -t nat -D" + " " + cmd)
	}
}

func ensureNewChain(cmd string) {
	out := RunCommand("iptables -w -t nat -N" + " " + cmd)
	if out == "iptables: Chain already exists." {
		fmt.Println("chain already exists, action will be canceled")
	}

}

var cnt = 0

func getFreeClusterIP() string {
	cnt++
	return fmt.Sprintf("10.111.111.%d", cnt)

}

func KubeProxyInit() {
	// 在nat表中创建规则链 KUBE-SERVICES
	ensureNewChain("KUBE-SERVICES -m comment --comment \"k8s service chain\"")
	// RunCommand("iptables -t nat -N KUBE-SERVICES -m comment --comment \"k8s service chain\"")
	// 把本机发出的流量劫持（bushi）到自建的规则链
	ensureRuleAppend("OUTPUT -j KUBE-SERVICES")
	// RunCommand("iptables -t nat -A OUTPUT -j KUBE-SERVICES")

	// RunCommand("iptables -N KUBE-MARK-MASQ")
	// RunCommand("iptables -N KUBE-POSTROUTING")

}

func RegisterService(service object.Service, endpoint object.Endpoint) {
	fmt.Printf("KubeProxy: register service : %s  %s\n", service.Metadata.Name, service.Metadata.Uid)
	fmt.Println(endpoint.PodIps)
	ports := service.Spec.Ports
	clusterIP := service.Spec.ClusterIP
	if clusterIP == "" {
		clusterIP = getFreeClusterIP()
	}

	for _, port := range ports {
		protocol := string(port.Protocol)
		chainId := service.Metadata.Uid[0:10]
		// 构造service该端口对应的SVC链名，chainId取service的uid的前十位
		// 创建链
		svcChain := fmt.Sprintf("KUBE-SVC-%s-%d", chainId, port.Port)
		// cmd := fmt.Sprintf("iptables -t nat -N %s -m comment --comment \"svc chain for service: %s port %d\"", svcChain, service.Metadata.Name, port.Port)
		cmd := fmt.Sprintf("%s -m comment --comment \"svc chain for service: %s port %d\"", svcChain, service.Metadata.Name, port.Port)
		ensureNewChain(cmd)
		// cmd = fmt.Sprintf("iptables -A KUBE-SERVICES -p %s -d %s/32 --dport %d -j KUBE-MARK-MASQ", protocol, clusterIP, port.Port)
		// 在KUBE-SERVICES链中增加到SVC链的转发规则（第一次转发，用以分开service各端口）
		// cmd = fmt.Sprintf("iptables -t nat -A KUBE-SERVICES -p %s -d %s/32 --dport %d -j %s", protocol, clusterIP, port.Port, svcChain)
		cmd = fmt.Sprintf("KUBE-SERVICES -p %s -d %s/32 --dport %d -j %s", protocol, clusterIP, port.Port, svcChain)
		ensureRuleAppend(cmd)

		podsLen := len(endpoint.PodIps)
		fmt.Println(podsLen)
		i := 0
		for _, podIp := range endpoint.PodIps {
			if podIp == "" {
				continue
			}
			// 构造对应service-port-pod的规则链（SEP）
			sepChain := fmt.Sprintf("KUBE-SEP-%s-%d-%d", chainId, port.Port, i)
			cmd := fmt.Sprintf("%s -m comment --comment \"sep chain %d for service: %s port %d\" ", sepChain, i, service.Metadata.Name, port.Port)
			ensureNewChain(cmd)
			// 在SVC链中增加跳转到SEP链的规则（第二次转发，用以在各pod间负载均衡，随机策略）
			if i == podsLen-1 {
				// cmd = fmt.Sprintf("iptables -t nat -A %s -j %s", svcChain, sepChain)
				cmd = fmt.Sprintf("%s -j %s", svcChain, sepChain)
			} else {
				pro := 1.0 / (float64(podsLen) - float64(i))
				// cmd = fmt.Sprintf("iptables -t nat -A %s -m statistic --mode random --probability %f -j %s", svcChain, pro, sepChain)
				cmd = fmt.Sprintf("%s -m statistic --mode random --probability %f -j %s", svcChain, pro, sepChain)

			}
			ensureRuleAppend(cmd)
			// 在SEP链上增加跳转到pod的规则（第三次转发，通过DNAT将目的地转换为 podIP:podPort
			// cmd = fmt.Sprintf("iptables -t nat -A %s -p %s -m tcp -j DNAT --to-destination %s:%d", sepChain, protocol, podIp, port.TargetPort)
			cmd = fmt.Sprintf("%s -p %s -m tcp -j DNAT --to-destination %s:%d", sepChain, protocol, podIp, port.TargetPort)
			ensureRuleAppend(cmd)
			i++
		}
	}
}

// 注册的镜像操作
func DeleteService(service object.Service, endpoint object.Endpoint) {
	fmt.Printf("delete service %s %s\n", service.Metadata.Name, service.Metadata.Uid)
	ports := service.Spec.Ports
	clusterIP := service.Spec.ClusterIP
	if clusterIP == "" {
		fmt.Println("cluster IP is NULL")
		return
	}
	for _, port := range ports {
		protocol := string(port.Protocol)
		chainId := service.Metadata.Uid[0:10]
		svcChain := fmt.Sprintf("KUBE-SVC-%s-%d", chainId, port.Port)
		// cmd := fmt.Sprintf("iptables -t nat -D KUBE-SERVICES -p %s -d %s/32 --dport %d -j %s", protocol, clusterIP, port.Port, svcChain)
		cmd := fmt.Sprintf("KUBE-SERVICES -p %s -d %s/32 --dport %d -j %s", protocol, clusterIP, port.Port, svcChain)
		ensureRuleDelete(cmd)
		RunCommand(cmd)
		RunCommand(fmt.Sprintf("iptables -t nat -F %s", svcChain))
		RunCommand(fmt.Sprintf("iptables -t nat -X %s", svcChain))
		podsLen := len(endpoint.PodIps)
		for i := 0; i < podsLen; i++ {
			sepChain := fmt.Sprintf("KUBE-SEP-%s-%d-%d", chainId, port.Port, i)
			RunCommand(fmt.Sprintf("iptables -t nat -F %s", sepChain))
			RunCommand(fmt.Sprintf("iptables -t nat -X %s", sepChain))
		}

	}
}

type KubeProxy struct {
	serviceSubscriber     *subscriber.Subscriber
	serviceQueue          string
	serviceHandler        serviceHandler
	EndpointSubscriberMap map[string]subscriber.Subscriber
}

// kubeproxy用于维护endpoint对象

func CreateKubeProxy() *KubeProxy {
	sub, _ := subscriber.NewSubscriber(global.MQHost)

	kubeProxy := KubeProxy{
		serviceSubscriber: sub,
		serviceQueue:      "services",
	}
	handler := serviceHandler{
		proxy: &kubeProxy,
	}
	kubeProxy.serviceHandler = handler
	kubeProxy.EndpointSubscriberMap = make(map[string]subscriber.Subscriber)
	return &kubeProxy
}

func (proxy *KubeProxy) Run() {
	KubeProxyInit()
	err := proxy.serviceSubscriber.Subscribe(proxy.serviceQueue, subscriber.Handler(proxy.serviceHandler))
	if err != nil {
		fmt.Printf(err.Error())
		_ = proxy.serviceSubscriber.CloseConnection()
	}
}

type serviceHandler struct {
	proxy *KubeProxy
}

func (h serviceHandler) Handle(jsonMsg []byte) {
	log.Println("Service get subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	service := object.Service{}
	prevService := object.Service{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &service)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevService)
	log.Println("type: " + string(rune(msg.EventType)))

	// 获取service对应的endpoint
	client := HTTPClient.CreateHTTPClient(global.ServerHost)
	getMsg, _ := json.Marshal(service.Metadata.Name)
	resp := client.Post("/endpoints/get", getMsg)

	fmt.Printf("Get endpoint: %s\n", resp)
	endpoint := object.Endpoint{}
	var epStr string
	json.Unmarshal([]byte(resp), &epStr)
	json.Unmarshal([]byte(epStr), &endpoint)
	fmt.Printf("len: %d\n", len(endpoint.PodIps))

	switch msg.EventType {
	// 创建service：配表，让该service监听endpoints队列
	case object.CREATE:
		RegisterService(service, endpoint)
		handler := endpointHandler{
			serviceName: service.Metadata.Name,
		}
		go func() {
			fmt.Printf("service: %s\n", service.Metadata.Uid)
			fmt.Printf("prev: %s\n", prevService.Metadata.Uid)
			sub, _ := subscriber.NewSubscriber(global.MQHost)
			h.proxy.EndpointSubscriberMap[service.Metadata.Name] = *sub
			err := sub.Subscribe("endpoints", subscriber.Handler(handler))
			if err != nil {
				fmt.Printf(err.Error())
				_ = sub.CloseConnection()
			}
		}()
	case object.DELETE:
		DeleteService(service, endpoint)
		sub := h.proxy.EndpointSubscriberMap[service.Metadata.Name]
		sub.CloseConnection()
		delete(h.proxy.EndpointSubscriberMap, service.Metadata.Name)

	}
}

type endpointHandler struct {
	serviceName string
}

func (h endpointHandler) Handle(jsonMsg []byte) {
	log.Println("endpoint get subscribe: " + string(jsonMsg))
	msg := object.MQMessage{}
	endpoint := object.Endpoint{}
	prevEndpoint := object.Endpoint{}
	_ = json.Unmarshal(jsonMsg, &msg)
	_ = json.Unmarshal([]byte(msg.Value), &endpoint)
	_ = json.Unmarshal([]byte(msg.PrevValue), &prevEndpoint)
	log.Println("type: " + string(rune(msg.EventType)))
	fmt.Printf("prevalue: %s\n", msg.PrevValue)
	fmt.Printf("value: %s\n", msg.Value)

	switch msg.EventType {
	case object.UPDATE:
		// 改动的endpoint属于当前service
		if endpoint.ServiceName == h.serviceName {
			// 获取当前service
			client := HTTPClient.CreateHTTPClient(global.ServerHost)
			getMsg, _ := json.Marshal(h.serviceName)
			resp := client.Post("/services/get", getMsg)
			fmt.Printf("Get service: %s\n", resp)
			service := object.Service{}
			var svStr string
			json.Unmarshal([]byte(resp), &svStr)
			json.Unmarshal([]byte(svStr), &service)

			currLen := len(endpoint.PodIps)
			prevLen := len(prevEndpoint.PodIps)
			// 改动的是pod数量，重新为service配表
			if currLen != prevLen {
				// pod删除
				// if currLen < prevLen {
				// 	for k, _ := range prevEndpoint.PodIps {
				// 		_, ok := endpoint.PodIps[k]
				// 		// 找到被删除的pod ip
				// 		if !ok {

				// 		}
				// 	}
				// }
				DeleteService(service, prevEndpoint)
				RegisterService(service, endpoint)
			} else {
				// 检查podip列表有无变化
				flag := false
				for k, v := range endpoint.PodIps {
					if prevEndpoint.PodIps[k] == "" && v != "" {
						flag = true
					}
				}
				if flag {
					DeleteService(service, prevEndpoint)
					RegisterService(service, endpoint)
				}
			}

		}
	}

}
