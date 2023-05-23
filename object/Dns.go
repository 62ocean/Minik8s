package object

type Path struct {
	Path        string
	ServiceName string `yaml:"serviceName"`
	ServicePort int    `yaml:"servicePort"`
}
type Host struct {
	HostName string `yaml:"hostName"`
	Paths    []Path
}
type DnsSpec struct {
	Hosts []Host
}
type Dns struct {
	Kind     string
	Metadata Metadata
	Spec     DnsSpec
}
