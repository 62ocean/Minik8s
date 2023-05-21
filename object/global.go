package object

import "time"

type Status int
type EventType int

const (
	RUNNING Status = 0
	STOPPED Status = 1
	PENDING Status = 2
)

<<<<<<< HEAD
type Metadata struct {
	Name   string `yaml:"name"`
	Labels Labels `yaml:"labels"`
	//Labels            map[string]string `yaml:"labels"`
	Namespace         string
	Uid               string
	CreationTimestamp time.Time
}

type Labels struct {
	App string `yaml:"app"`
	Env string `yaml:"env"`
=======
const (
	CREATE EventType = 0
	UPDATE EventType = 1
	DELETE EventType = 2
)

type MQMessage struct {
	EventType EventType
	Value     string
>>>>>>> fcdca8038c800a8c1d95a2db516d70aff93a02b9
}
