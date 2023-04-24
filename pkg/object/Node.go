package object

type Node struct {
	Name   string
	IP     string
	Status NodeStatus
}

type NodeStatus int

const (
	RUN   NodeStatus = 0
	STOP  NodeStatus = 1
	ERROR NodeStatus = 2
)
