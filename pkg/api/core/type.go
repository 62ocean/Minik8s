package core

import (
	"time"
)

// 所有API对象的父类，其属性是API对象的共有属性
type ObjectMeta struct {
	Name              string
	Namespace         string
	UUID              string
	CreationTimestamp time.Time
	Labels            map[string]string
}

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

type Volume struct {
	Name string
	//VolumeSource
}

type ContainerPort struct {
	Name          string
	Protocol      Protocol
	ContainerPort int32
}

type VolumeMount struct {
	Name      string
	MountPath string
}
type Resource struct {
	CPU    int64
	Memory int64
}
type Container struct {
	Name         string
	Image        string
	Command      []string
	Args         []string
	Ports        []ContainerPort
	VolumeMounts []VolumeMount
	Resource
	//Env          []EnvVar
}
type PodSpec struct {
	Volumes    []Volume
	Containers []Container
}

type PodStatus struct {
	Phase PodPhase
}
type Pod struct {
	Kind       string
	ObjectMeta `yaml：“metadata”`
	Spec       PodSpec
	Status     PodStatus
}

type ServicePort struct {
	Name     string
	Protocol Protocol
	// 对外暴露的端口
	Port     int32
	NodePort int32
	// 对pod暴露的端口
	TargetPort int32
}
type ServiceType string

type ServiceSpec struct {
	Type      ServiceType
	Ports     []ServicePort
	Selector  map[string]string
	ClusterIP string
	pods      []Pod
}

type Service struct {
	Kind       string
	ObjectMeta `yaml：“metadata”`
	Spec       ServiceSpec
}
