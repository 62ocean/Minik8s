package object

type GPUJob struct {
	ApiVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   Metadata   `yaml:"metadata"`
	Spec       GPUJobSpec `yaml:"spec"`
	Status     Status
	Output     string
	Error      string
}

type GPUJobSpec struct {
	Nodes           int      `yaml:"nodes"`
	NumTasksPerNode int      `yaml:"numTasksPerNode"`
	CpusPerTask     int      `yaml:"cpusPerTask"`
	NumGpus         int      `yaml:"numGpus"`
	CompileScripts  []string `yaml:"compileScripts"`
	Program         string   `yaml:"program"`
	Executable      string   `yaml:"exe"`
}

// ScriptTemplate 使用不转义字符串，所见即所得
const ScriptTemplate = `#!/bin/bash
#SBATCH --job-name=%s
#SBATCH --partition=dgx2
#SBATCH --nodes=%d
#SBATCH --ntasks-per-node=%d
#SBATCH --cpus-per-task=%d
#SBATCH --gres=gpu:%d
#SBATCH --mail-type=end
#SBATCH --mail-user=rongchuan_liu@sjtu.edu.cn  
#SBATCH --output=%%j.out
#SBATCH --error=%%j.out

%s
`
