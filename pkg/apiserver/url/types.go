package url

type Node struct {
	Id   string
	Name string
}

type NodeResource struct {
	users map[string]Node
}
