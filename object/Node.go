package object

type NodeStorage struct {
	Node   Node
	Status Status
}

type Node struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   Metadata
	IP         string
}
