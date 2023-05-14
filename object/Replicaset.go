package object

type ReplicaSet struct {
	ApiVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       ReplicasetSpec `yaml:"spec"`
}

type ReplicasetSpec struct {
	Replicas    int         `yaml:"replicas"`
	Selector    Selector    `yaml:"selector"`
	PodTemplate PodTemplate `yaml:"template"`
}

type Selector struct {
	MatchLabels MatchLabels `yaml:"matchLabels"`
}

type MatchLabels struct {
	App string `yaml:"app"`
	Env string `yaml:"env"`
}

type PodTemplate struct {
	Metadata `yaml:"metadata"`
	PodSpec  PodSpec `yaml:"spec"`
}
