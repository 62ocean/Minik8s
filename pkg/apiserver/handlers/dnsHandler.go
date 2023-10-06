package handlers

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/emicklei/go-restful/v3"
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

func CreateDns(request *restful.Request, response *restful.Response) {
	fmt.Println("apiserver handler: create dns")

	dns := new(object.Dns)
	err := request.ReadEntity(&dns)
	if err != nil {
		log.Println(err)
		return
	}

	coreFile, err := os.OpenFile("build/Corefile", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	coreFile2, err2 := os.OpenFile("Corefile", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	check(err)
	check(err2)
	defer coreFile.Close()
	defer coreFile2.Close()

	dnsStr, _ := json.Marshal(*dns)
	etcd.Put("/registry/dns", string(dnsStr))
	fmt.Println(string(dnsStr))

	for i, host := range dns.Spec.Hosts {
		hostIp := fmt.Sprintf("%s.%d", global.HostNameIpPrefix, i)
		// 切分域名
		sep := "."
		arr := strings.Split(host.HostName, sep)
		// 根据域名构建存储键值
		key := "/coredns"
		for i := len(arr) - 1; i >= 0; i-- {
			key = fmt.Sprintf("%s/%s", key, arr[i])
		}
		val := fmt.Sprintf(" {\"host\":\"%s\",\"port\":80} ", hostIp)
		fmt.Println(key)
		fmt.Println(val)
		// 持久化到etcd
		etcd.Put(key, val)

		// 配置coreFile文件，没有就创建，以追加模式写入
		block := fmt.Sprintf(
			"%s {\n"+
				"  etcd {\n"+
				"    endpoint http://%s\n"+
				"    path /coredns\n"+
				"  }\n"+
				"}\n", host.HostName, global.EtcdHost)
		coreFile.WriteString(block)
		coreFile2.WriteString(block)
	}
	RunCommand("killall nginx")
	RunCommand("killall coredns")

	RunCommand("systemctl restart nignx")
	RunCommand("build/coredns &")

}

func check(err error) {

}

func GetDns(request *restful.Request, response *restful.Response) {
	dnsStr := etcd.GetOne("/registry/dns")
	msg, _ := json.Marshal(dnsStr)
	_, err := response.Write(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}
