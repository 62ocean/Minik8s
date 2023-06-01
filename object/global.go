package object

import "time"

type Status int
type EventType int

const (
	RUNNING  Status = 0
	STOPPED  Status = 1
	PENDING  Status = 2
	FINISHED Status = 3
)

const (
	CREATE EventType = 0
	UPDATE EventType = 1
	DELETE EventType = 2
)

type MQMessage struct {
	EventType EventType
	Value     string
	PrevValue string
}

type Metadata struct {
	Name              string `yaml:"name"`
	Labels            Labels `yaml:"labels"`
	Namespace         string
	Uid               string
	CreationTimestamp time.Time
}

type Labels struct {
	App string `yaml:"app"`
	Env string `yaml:"env"`
}

func (s Status)ToString() string{
	switch s{
	case 0:
		return "RUNNING"
	case 1:
		return "STOPPED"
	case 2:
		return "PENDING"
	case 3:
		return "FINISHED"
	default:
		return ""
	}
}
