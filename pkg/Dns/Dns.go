package Dns

import (
	"fmt"
	"k8s/object"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadDnsConfig(path string) (object.Dns, error) {
	dns := object.Dns{}
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return dns, err
	}

	err = yaml.Unmarshal(dataBytes, &dns)
	if err != nil {
		fmt.Println("解析yaml文件失败:", err)
		return dns, err
	}
	return dns, nil
}

func ConfigTest() {
	fmt.Println("测试 Dns 配置文件读取")
	dns, err := ReadDnsConfig("pkg/Dns/dnsConfigTest.yaml")
	if err != nil {
		fmt.Println("测试失败")
		return
	}
	fmt.Printf("解析结果：\n + dns -> %+v\n", dns)
	return
}

func CreateDns() {
	// dns, err := ReadDnsConfig(filePath)
	// if err != nil {
	// 	fmt.Println("CreateDns: read error")
	// 	return
	// }
	coreFile, err := os.OpenFile("./Corefile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	check(err)
	defer coreFile.Close()
	// nginxConfig, err := os.OpenFile("/etc/nginx/conf.d/default.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	// check(err)
	// defer nginxConfig.Close()

	// for i, host := range dns.Spec.Hosts {
	// 	hostIp := fmt.Sprintf("%s.%d", global.HostNameIpPrefix, i)
	// 	// 切分域名
	// 	sep := "."
	// 	arr := strings.Split(host.HostName, sep)
	// 	// 根据域名构建存储键值
	// 	key := "/coredns"
	// 	for i := len(arr) - 1; i >= 0; i-- {
	// 		key = fmt.Sprintf("%s/%s", key, arr[i])
	// 	}
	// 	val := fmt.Sprintf(" {\"host\":\"%s\",\"port\":80} ", hostIp)
	// 	fmt.Println(key)
	// 	fmt.Println(val)
	// 	// 持久化到etcd
	// 	etcd.Put(key, val)
	// 	// 配置coreFile文件，没有就创建，每次写入清空先前内容
	// 	block := fmt.Sprintf(
	// 		"%s {\n"+
	// 			"  etcd {\n"+
	// 			"    endpoint http://%s\n"+
	// 			"    path /coredns\n"+
	// 			"  }\n"+
	// 			"}\n", host.HostName, global.EtcdHost)
	// 	coreFile.WriteString(block)

	// 	// 依据路径配置nginx
	// 	nginxConfig.WriteString("server {\n")
	// 	nginxBlock := fmt.Sprintf("  listen 80;\n  server_name %s;\n", hostIp)
	// 	nginxConfig.WriteString(nginxBlock)
	// 	for _, path := range host.Paths {
	// 		service := object.Service{}
	// 		str := etcd.GetOne("/registry/services/" + path.ServiceName)
	// 		json.Unmarshal([]byte(str), &service)
	// 		nginxBlock = fmt.Sprintf("  location %s {\n    proxy_pass http://%s:%d/;\n  }\n",
	// 			path.Path, service.Spec.ClusterIP, path.ServicePort)
	// 		nginxConfig.WriteString(nginxBlock)

	// 	}
	// 	nginxConfig.WriteString("}")
	// }
	// 其余流量转发到DNS服务器
	block := fmt.Sprintf(". {\n  forward . 114.114.114.114\n  cache 30\n}\n")
	coreFile.WriteString(block)

}

func check(err error) {

}
