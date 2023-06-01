package object

type Hpa struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       HpaSpec  `yaml:"spec"`
}

type HpaSpec struct {
	ScaleTargetRef ScaleTargetRef `yaml:"scaleTargetRef"`
	MinReplicas    int            `yaml:"minReplicas"`
	MaxReplicas    int            `yaml:"maxReplicas"`
	Metrics        []Metric       `yaml:"metrics"`
}

type ScaleTargetRef struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
}

type Metric struct {
	Type     string   `yaml:"type"`
	Resource Resource `yaml:"resource"`
}

type Resource struct {
	Name   string `yaml:"name"`
	Target Target `yaml:"target"`
}

type Target struct {
	Type               string  `yaml:"type"`
	AverageUtilization float64 `yaml:"averageUtilization"`
}
