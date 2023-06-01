package pod

import (
	"fmt"
	"os/exec"
)

type VxlanDevice struct {
	// 设备名
	Name string
	// 设备id
	Vni     int
	DstPort int
	// 最大传输单元
	Mtu          int
	MasterBridge string
}

type NodeNetwork struct {
	Docker0MacAddr string
	IpAddr         string
	Subnet         string
	Gateway        string
	SubnetPrefix   string
}

func NewVxlanDevice(Name string, Vni int, DstPort int, MasterBridge string) *VxlanDevice {
	return &VxlanDevice{Name: Name, Vni: Vni, DstPort: DstPort, Mtu: 1450, MasterBridge: MasterBridge}
}

func RunCommand(cmd string) {
	fmt.Printf("RunCmd: %s\n", cmd)
	command := exec.Command("/bin/bash", "-c", cmd)
	if _, err := command.CombinedOutput(); err != nil {
		// panic("ERROR: " + err.Error())
	}
}

func (vx *VxlanDevice) Create() {
	fmt.Printf("create vxlan device...\n")
	// 添加vxlan设备
	cmd := fmt.Sprintf("ip link add %s type vxlan id %d dstport %d", vx.Name, vx.Vni, vx.DstPort)
	RunCommand(cmd)
	cmd = fmt.Sprintf("ip link set dev %s mtu %d", vx.Name, vx.Mtu)
	RunCommand(cmd)
	// 启用设备
	cmd = fmt.Sprintf("ip link set dev %s up", vx.Name)
	RunCommand(cmd)
	// 把vxlan设备连到MasterBridge（docker0）上
	cmd = fmt.Sprintf("ip link set dev %s master %s", vx.Name, vx.MasterBridge)
	RunCommand(cmd)
	fmt.Printf("The vxlan device has successfully created\n")
}

// 有其他设备加入当前网络，配置本节点的路由和fdb和arp信息
func (vx *VxlanDevice) AddNodeToNetwork(subnet string, gateway string, docker0MacAddr string, ipAddr string) {
	fmt.Printf("Add the new Node to the flannel network：subnet: %s, gateway: %s, docker0MacAddr: %s, ipAddr: %s\n", subnet, gateway, docker0MacAddr, ipAddr)
	cmd := fmt.Sprintf("ip route add %s via %s dev %s onlink", subnet, gateway, vx.Name)
	RunCommand(cmd)
	cmd = fmt.Sprintf("ip nei add %s dev %s lladdr %s", gateway, vx.Name, docker0MacAddr)
	RunCommand(cmd)
	cmd = fmt.Sprintf("bridge fdb add %s dev %s dst %s", docker0MacAddr, vx.Name, ipAddr)
	RunCommand(cmd)
}

func (vx *VxlanDevice) DelNodeFromNetwork(subnet string, gateway string, docker0MacAddr string, ipAddr string) {
	fmt.Printf("Delete the Node from the flannel network：subnet: %s, gateway: %s, docker0MacAddr: %s, ipAddr: %s\n", subnet, gateway, docker0MacAddr, ipAddr)
	cmd := fmt.Sprintf("ip route del %s via %s dev %s onlink", subnet, gateway, vx.Name)
	RunCommand(cmd)
	cmd = fmt.Sprintf("ip nei del %s dev %s lladdr %s", gateway, vx.Name, docker0MacAddr)
	RunCommand(cmd)
	cmd = fmt.Sprintf("bridge fdb del %s dev %s dst %s", docker0MacAddr, vx.Name, ipAddr)
	RunCommand(cmd)
}
