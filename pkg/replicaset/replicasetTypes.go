package replicaset

import "k8s/pkg/api/pod"

type ReplicaSet struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Spec struct {
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
	PodMetadata pod.Metadata `yaml:"metadata"`
	PodSpec     pod.Spec     `yaml:"spec"`
}
