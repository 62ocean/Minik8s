package GPU

import (
	"fmt"
	"os"
	"path"
)

type GPUServer struct {
	Client *SSHClient
}

// NewServer GPUServer对象的构造函数
func NewServer(name string) *GPUServer {
	server := GPUServer{
		Client: NewClient(),
	}
	return &server
}

// Run kubelet运行的入口函数
func (s *GPUServer) Run() {

}

func (s *GPUServer) submitJob() (err error) {
	if s.jobID, err = s.cli.SubmitJob(s.scriptPath()); err == nil {
		fmt.Printf("submit succeed, got jod ID: %s\n", s.jobID)
	}
	return err
}

func (s *GPUServer) prepare() (err error) {
	cudaFiles := s.getCudaFiles()
	if len(cudaFiles) == 0 {
		return fmt.Errorf("no available cuda files")
	}
	if err = s.uploadSmallFiles(cudaFiles); err != nil {
		return err
	}
	fmt.Println("upload cuda files successfully")
	if err = s.compile(); err != nil {
		return err
	}
	fmt.Println("compile successfully")
	if err = s.createJobScript(); err != nil {
		return err
	}
	fmt.Println("create job script successfully")
	return nil
}

func (s *GPUServer) downloadResult() {
	outputFile := s.args.Output
	if content, err := s.cli.ReadFile(outputFile); err == nil {
		if file, err := os.Create(path.Join(s.jobsURL, outputFile)); err == nil {
			defer file.Close()
			_, _ = file.Write([]byte(content))
		}
	}

	errorFile := s.args.Error
	if content, err := s.cli.ReadFile(errorFile); err == nil {
		if file, err := os.Create(path.Join(s.jobsURL, errorFile)); err == nil {
			defer file.Close()
			_, _ = file.Write([]byte(content))
		}
	}
}
