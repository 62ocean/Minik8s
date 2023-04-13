package pod

type ContainerPort struct {
	Port int32 `yaml:"containerPort"`
}

type Container struct {
	Name  string          `yaml:"name"`
	Image string          `yaml:"image"`
	Ports []ContainerPort `yaml:"ports"`
}

type Spec struct {
	Containers []Container `yaml:"containers"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Pod struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
