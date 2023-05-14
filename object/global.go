package object

type Status int

const (
	RUNNING Status = 0
	STOPPED Status = 1
	PENDING Status = 2
)

type Metadata struct {
	Name      string `yaml:"name"`
	Labels    Labels `yaml:"labels"`
	Namespace string
	Uid       string
}

type Labels struct {
	App string `yaml:"app"`
	Env string `yaml:"env"`
}
