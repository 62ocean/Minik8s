package flannel

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"k8s/pkg/api/pod"
	"k8s/pkg/etcd"
	"k8s/pkg/global"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Config struct {
	etcdEndpoint string
	etcdPrefix   string
	dockerNet    string
	//主网段
	flannelNetwork string
	//子网段
	flannelSubnet        string
	flannelSubnetPrefix  string
	flannelSubnetGateway string
	nodeIP               string
	nodeID               string
}

var conf Config

const vni = 1
const vxlanDstPort = 8472
const flannelNetworkPrefix = "162.16"
const flannelNetwork = "162.16.0.0/16"
const etcdPrefix = "/flannel.com/network"

func readConfiguration(configurationFile string) map[string]string {
	var properties = make(map[string]string)
	confFile, err := os.OpenFile(configurationFile, os.O_RDONLY, 0666)
	defer func(confFile *os.File) {
		if err := confFile.Close(); err != nil {
			panic(err)
		}
	}(confFile)
	if err != nil {
		fmt.Printf("The config file %s is not exits.", configurationFile)
	} else {
		reader := bufio.NewReader(confFile)
		for {
			if confString, err := reader.ReadString('\n'); err != nil {
				if err == io.EOF {
					break
				}
			} else {
				if len(confString) == 0 || confString == "\n" || confString[0] == '#' {
					continue
				}
				properties[strings.Split(confString, "=")[0]] = strings.Replace(strings.Split(confString, "=")[1], "\n", "", -1)
			}
		}
	}
	return properties
}
func ConfigInit() {
	//if dataBytes, err := os.ReadFile("pkg/apiserver/flannel/flannel.properties"); err != nil {
	//	panic("读取文件失败：" + err.Error())
	//} else {
	//	fmt.Printf("文件：%s", dataBytes)
	//	conf = Config{}
	//	err = yaml.Unmarshal(dataBytes, &conf)
	//	fmt.Printf("解析结果：\n + service -> %+v\n", conf)
	//	if err != nil {
	//		fmt.Printf("yaml 解析失败")
	//	}
	//}

	// config := readConfiguration("pkg/apiserver/flannel/flannel.properties")
	config := readConfiguration("build/flannel.properties")

	conf.etcdEndpoint = config["etcd-endpoint"]
	conf.etcdPrefix = etcdPrefix
	conf.nodeIP = config["node-ip"]
	conf.nodeID = config["node-id"]
	conf.dockerNet = fmt.Sprintf("%s.%s.1/24", flannelNetworkPrefix, conf.nodeID)
	conf.flannelNetwork = flannelNetwork
	conf.flannelSubnet = fmt.Sprintf("%s.%s.0/24", flannelNetworkPrefix, conf.nodeID)
	conf.flannelSubnetPrefix = fmt.Sprintf("%s.%s", flannelNetworkPrefix, conf.nodeID)
	conf.flannelSubnetGateway = fmt.Sprintf("%s.%s.1", flannelNetworkPrefix, conf.nodeID)

	fmt.Printf("==================================== Starting Flanneld ===================================\n")
	fmt.Printf("        etcd-endpoint:%s\n", conf.etcdEndpoint)
	fmt.Printf("           etcd-prefix:%s\n", conf.etcdPrefix)
	fmt.Printf("             docker-net:%s\n", conf.dockerNet)
	fmt.Printf("       flannel-network:%s\n", conf.flannelNetwork)
	fmt.Printf("        flannel-subnet:%s\n", conf.flannelSubnet)
	fmt.Printf("flannel-subnet-gateway:%s\n", conf.flannelSubnetGateway)
	fmt.Printf("             node-ip:%s\n", conf.nodeIP)
	fmt.Printf("             node-id:%s\n", conf.nodeID)
	fmt.Printf("=========================================================================================\n")
}

func SetDockerBipNet() {
	// daemonJsonPath := `"{
	// \"bip\":\"` + conf.dockerNet + `\"}"`
	// cmd := "echo " + daemonJsonPath + " > /etc/docker/daemon.json"

	config, _ := os.OpenFile("/etc/docker/daemon.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	config.WriteString("{\n")

	cmd := fmt.Sprintf("  \"bip\":\"%s\",\n", conf.dockerNet)
	config.WriteString(cmd)
	cmd = fmt.Sprintf("  \"dns\":[\"%s\"]\n", global.NameServerIp)
	config.WriteString(cmd)
	config.WriteString("}\n")
	config.Close()

	pod.RunCommand(cmd)
	pod.RunCommand("systemctl restart docker")
	fmt.Printf("finish SetDockerBipNet\n")
}

func InitIptables() {
	// 设置filter表FORWARD链默认规则为允许通过
	pod.RunCommand("iptables -t filter -P FORWARD ACCEPT")
	// nat表POSTROUTING链插入规则，允许目标为flannelNetwork的包通过
	cmd := fmt.Sprintf("iptables -t nat -I POSTROUTING -d %s -j ACCEPT", conf.flannelNetwork)
	pod.RunCommand(cmd)
}

func GetLocalMacAddr(devName string) string {
	fmt.Println(devName)
	interfaces, _ := net.InterfaceByName(devName)
	macAddr := fmt.Sprintf("%v", interfaces.HardwareAddr)
	fmt.Printf("GetLocalMacAddr: %s\n", macAddr)
	return macAddr
}

func GetLocalNodeNetwork() pod.NodeNetwork {
	return pod.NodeNetwork{
		IpAddr:         conf.nodeIP,
		Docker0MacAddr: GetLocalMacAddr("docker0"),
		Subnet:         conf.flannelSubnet,
		Gateway:        conf.flannelSubnetGateway,
		SubnetPrefix:   conf.flannelSubnetPrefix,
	}
}

// 对etcd中已存在的节点配置本地路由等
func addCurrentNodes(vx *pod.VxlanDevice) {
	currentNodes := etcd.GetDirectory(conf.etcdPrefix)
	for _, v := range currentNodes {
		node := pod.NodeNetwork{}
		_ = json.Unmarshal([]byte(v), &node)
		vx.AddNodeToNetwork(node.Subnet, node.Gateway, node.Docker0MacAddr, node.IpAddr)
	}
}

func SetupCloseHandler(vx *pod.VxlanDevice) {
	c := make(chan os.Signal, 0)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(vx *pod.VxlanDevice) {
		<-c
		fmt.Printf("Interrupt\n")
		cmd := fmt.Sprintf("ip link del %s", vx.Name)
		pod.RunCommand(cmd)
		etcd.Del(conf.etcdPrefix + "/" + conf.nodeIP)
		fmt.Printf("Delete vxlan dev\n")
		os.Exit(0)
	}(vx)

}

func main() {

}

func Exec() {

	ConfigInit()
	SetDockerBipNet()

	InitIptables()
	vx := pod.NewVxlanDevice("vxlan0", vni, vxlanDstPort, "docker0")
	val, _ := json.Marshal(GetLocalNodeNetwork())
	cli := etcd.GetEtcdClient(conf.etcdEndpoint)
	if cli == nil {
		fmt.Printf("connect failed\n")
	}
	vx.Create()
	SetupCloseHandler(vx)
	// 为flannel网络当前存在的节点配置路由等
	addCurrentNodes(vx)
	// 将当前节点注册到etcd
	etcd.Put(conf.etcdPrefix+"/"+conf.nodeIP, string(val))
	etcd.WatchPrefix(conf.etcdPrefix, vx, conf.nodeIP, cli)
}
