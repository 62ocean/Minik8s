package object

type ContainerPort struct {
	Port int32 `yaml:"containerPort"`
}

type ContainerVolumeMountConfig struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
	ReadOnly  bool   `yaml:"readOnly"`
}

type ContainerResources struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type ContainerResourcesConfig struct {
	Limits ContainerResources `yaml:"limits"`
}

type Container struct {
	Name         string                       `yaml:"name"`
	Image        string                       `yaml:"image"`
	Ports        []ContainerPort              `yaml:"ports"`
	Command      []string                     `yaml:"command"`
	Args         []string                     `yaml:"args"`
	Resources    ContainerResourcesConfig     `yaml:"resources"`
	VolumeMounts []ContainerVolumeMountConfig `yaml:"volumeMounts"`
	CopyFile     string
	CopyDst      string
}

type VolumeEmptyDirConfig struct {
}

type VolumeHostPathConfig struct {
	Path string `yaml:"path"`
}

type VolumeConfig struct {
	Name     string               `yaml:"name"`
	EmptyDir VolumeEmptyDirConfig `yaml:"emptyDir"`
	HostPath VolumeHostPathConfig `yaml:"hostPath"`
}

type PodSpec struct {
	Containers []Container    `yaml:"containers"`
	Volumes    []VolumeConfig `yaml:"volumes"`
}

type Pod struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	IP         string
	Metadata   Metadata `yaml:"metadata"`
	Spec       PodSpec  `yaml:"spec"`
}

type PodStorage struct {
	Config         Pod
	Status         Status
	Node           string //放node的uid
	Replica        int
	RunningMetrics RunningMetrics
}

type RunningMetrics struct {
	CPUUtil float64
	MemUtil float64
}
