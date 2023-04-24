package object

type NodeStorage struct {
	Node   Node
	Status Status
}

type Node struct {
	Name string
	IP   string
}
