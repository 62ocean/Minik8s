package object

type NodeStorage struct {
	Node   Node
	Status Status
}

type Node struct {
	Metadata Metadata
	IP       string
}
