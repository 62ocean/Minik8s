package object

// 所有API对象的父类，其属性是API对象的共有属性
//type ObjectMeta struct {
//	Name              string
//	Namespace         string
//	UUID              string
//	CreationTimestamp time.Time
//	Labels            map[string]string
//}

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

//type Volume struct {
//	Name string
//	//VolumeSource
//}

//type ContainerPort struct {
//	Name          string
//	Protocol      Protocol
//	ContainerPort int32
//}

//type VolumeMount struct {
//	Name      string
//	MountPath string
//}
//type Resource struct {
//	CPU    int64
//	Memory int64
//}
//type Container struct {
//	Name         string
//	Image        string
//	Command      []string
//	Args         []string
//	Ports        []ContainerPort
//	VolumeMounts []VolumeMount
//	Resource
//	//Env          []EnvVar
//}
//type PodSpec struct {
//	Volumes    []Volume
//	Containers []Container
//}

//type PodStatus struct {
//	Phase PodPhase
//}
//type Pod struct {
//	Kind       string
//	ObjectMeta `json:“metadata”`
//	Spec       PodSpec
//	Status     PodStatus
//}

type ServicePort struct {
	Name     string
	Protocol Protocol
	// port是k8s集群内部访问service的端口
	Port int32
	// 外部访问k8s集群中service的端口
	NodePort int32
	// pod的端口
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
	Kind     string
	Metadata Metadata
	Spec     ServiceSpec
}
