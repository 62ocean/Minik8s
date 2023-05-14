package object

type Status int
type EventType int

const (
	RUNNING Status = 0
	STOPPED Status = 1
	PENDING Status = 2
)

const (
	CREATE EventType = 0
	UPDATE EventType = 1
	DELETE EventType = 2
)

type MQMessage struct {
	EventType EventType
	Value     string
}

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