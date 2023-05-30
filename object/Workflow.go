package object

type Workflow struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Start      string   `yaml:"start"`
	Params     []Param  `yaml:"params"`
	Steps      []Step   `yaml:"steps"`
}

type Param struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

type Step struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Next string `yaml:"next"`

	Choices []Choice `yaml:"choices"`
}

type Choice struct {
	Variable string `yaml:"variable"`
	Type     string `yaml:"type"`
	Value    int    `yaml:"value"`
	Next     string `yaml:"next"`
}
