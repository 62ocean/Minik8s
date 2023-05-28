package object

type Endpoint struct {
	ServiceName string
	Selector    Labels
	// pod-id 到 pod-ip
	PodIps map[string]string
}
