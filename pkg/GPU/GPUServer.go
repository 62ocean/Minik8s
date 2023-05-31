package GPU

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"time"
)

type GPUServer struct {
	Client  *SSHClient
	Job     object.GPUJob
	HTTPCli *HTTPClient.Client
}

// NewServer GPUServer对象的构造函数
func NewServer(job object.GPUJob) *GPUServer {
	// 建立ssh连接
	server := GPUServer{
		Client:  NewClient(),
		Job:     job,
		HTTPCli: HTTPClient.CreateHTTPClient(global.ServerHost),
	}
	return &server
}

// Run GPUServer
func (s *GPUServer) Run() {
	// 上传文件
	s.uploadFile()
	//生成slurm脚本
	s.generateScript()
	// 提交作业
	s.submitJob()
	// 获取结果
	for {
		time.Sleep(time.Second * 5)
		if !s.isJobRunning() {
			break
		}
	}
	s.getResult()
	s.Client.Close()
}

func (s *GPUServer) uploadFile() {
	fmt.Println("Begin To Upload File")
	dstDir := "/lustre/home/acct-stu/stu1639/" + s.Job.Metadata.Name + "/"
	if resp, err := s.Client.Mkdir(dstDir); err != nil {
		fmt.Println(resp)
	}
	s.Client.WorkDir = dstDir
	err := s.Client.Cli.Upload(s.Job.Spec.Program, dstDir+s.Job.Spec.Program)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *GPUServer) generateScript() {
	fmt.Println("Begin To Generate Script")
	var runCMD string
	// 编译命令也加进去
	for _, cmd := range s.Job.Spec.CompileScripts {
		runCMD += cmd + "\n"
	}
	runCMD += "./" + s.Job.Spec.Executable
	script := fmt.Sprintf(
		object.ScriptTemplate,
		s.Job.Metadata.Name,
		s.Job.Spec.Nodes,
		s.Job.Spec.NumTasksPerNode,
		s.Job.Spec.CpusPerTask,
		s.Job.Spec.NumGpus,
		runCMD,
	)
	scriptName := s.Job.Metadata.Name + ".slurm"
	resp, err := s.Client.WriteFile(scriptName, script)
	if err != nil {
		fmt.Println(resp)
	}
}

func (s *GPUServer) submitJob() {
	scriptName := s.Job.Metadata.Name + ".slurm"
	if resp, err := s.Client.SubmitJob(scriptName); err == nil {
		fmt.Printf("get rsponse after submit: %s\n", resp)
		s.Job.Metadata.Uid = resp[len(resp)-9 : len(resp)-1]
		fmt.Printf("submit succeed, got jod ID: %s\n", s.Job.Metadata.Uid)
		s.updateStatus(object.RUNNING)
	}
}

func (s *GPUServer) isJobRunning() bool {
	return s.Client.isJobRunning(s.Job.Metadata.Uid)
}

func (s *GPUServer) getResult() {
	s.Job.Output = s.Client.getOutPut(s.Job.Metadata.Uid)
	s.Job.Error = s.Client.getError(s.Job.Metadata.Uid)
	fmt.Println("Job Output: " + s.Job.Output)
	s.updateStatus(object.FINISHED)
	s.deletePod()
}

func (s *GPUServer) updateStatus(status object.Status) {
	s.Job.Status = status
	jobInfo, _ := json.Marshal(s.Job)
	s.HTTPCli.Post("/gpuJobs/update", jobInfo)
}

func (s *GPUServer) deletePod() {
	podName := "GPUJob_" + s.Job.Metadata.Name
	name, _ := json.Marshal(podName)
	s.HTTPCli.Post("/pods/remove", name)
}
