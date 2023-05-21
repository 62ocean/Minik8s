package kube_proxy

import (
	"bytes"
	"fmt"
	"k8s/object"
	"os/exec"
)

func RunCommand(cmd string) {
	fmt.Printf("RunCmd: %s\n", cmd)
	command := exec.Command("/bin/bash", "-c", cmd)
	if _, err := command.CombinedOutput(); err != nil {
		panic("ERROR: " + err.Error())
	}
}

var cnt = 0

func getFreeClusterIP() string {
	cnt++
	return fmt.Sprintf("10.111.111.%d", cnt)

}

func serviceInit() {
	// 创建规则链 KUBE-SERVICES
	cmd := fmt.Sprintf("iptables -N KUBE-SERVICES")
	RunCommand(cmd)
}

func registerService(service object.Service) {
	ports := service.Spec.Ports
	clusterIP := service.Spec.ClusterIP
	if clusterIP == "" {
		clusterIP = getFreeClusterIP()
	}

	for _, port := range ports {
		protocol := string(port.Protocol)
		svcChain := fmt.Sprintf("KUBE-SVC-%s%d", bytes.ToUpper([]byte(service.Metadata.Name)), port.Port)
		cmd := fmt.Sprintf("iptables -N %s", svcChain)
		RunCommand(cmd)
		cmd = fmt.Sprintf("iptables -A KUBE-SERVICES -p %s -d %s/32 --dport %d -j %s", protocol, clusterIP, port.Port, svcChain)
		RunCommand(cmd)

		pods := service.Spec.Pods
		podsLen := len(service.Spec.Pods)
		for i, pod := range pods {
			sepChain := fmt.Sprintf("KUBE-SEP-%s-POD%d", bytes.ToUpper([]byte(service.Metadata.Name)), i)
			RunCommand(fmt.Sprintf("iptables -N %s", sepChain))
			if i == podsLen-1 {
				cmd = fmt.Sprintf("iptables -A %s -j %s", svcChain, sepChain)
			} else {
				pro := 1.0 / (float64(podsLen) - float64(i))
				cmd = fmt.Sprintf("iptables -A %s --probability %f -j %s", svcChain, pro, sepChain)
			}
			RunCommand(cmd)

			cmd = fmt.Sprintf("iptables -A %s -p %s -j DNAT --to-destination %s:%d", sepChain, protocol, pod.IP, port.TargetPort)
		}
	}
}
