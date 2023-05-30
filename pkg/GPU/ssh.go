package GPU

import (
	"fmt"
	"github.com/melbahja/goph"
	"strings"
)

const (
	gpuUser      = "stu1639"
	gpuPasswd    = "34j3oQ*%"
	gpuLoginAddr = "pilogin.hpc.sjtu.edu.cn"
	gpuDataAddr  = "data.hpc.sjtu.edu.cn"
	gpuUserDir   = "/lustre/home/acct-stu/stu1639"
)

type SSHClient struct {
	Username string
	Password string
	Cli      *goph.Client
	WorkDir  string
}

/*--------------------SSH CLIENT--------------------------*/

func NewClient() *SSHClient {
	sshCli := newSSHClient(gpuUser, gpuPasswd)
	return &SSHClient{
		Username: gpuUser,
		Password: gpuPasswd,
		Cli:      sshCli,
	}
}

func (client *SSHClient) Compile(scripts []string) (string, error) {
	var resp []byte
	var err error
	for _, cmd := range scripts {
		resp, err = client.Cli.Run(cmd)
	}
	return string(resp), err
}

func (client *SSHClient) SubmitJob(filename string) (string, error) {
	filePath := client.WorkDir + filename
	cmd := "sbatch " + filePath
	resp, err := client.Cli.Run(cmd)
	return string(resp), err
}

func (client *SSHClient) isJobRunning(jobID string) bool {
	res, _ := client.Cli.Run("squeue | grep " + jobID)
	fmt.Println("Get job status: " + string(res))
	return len(res) > 5
}

func (client *SSHClient) getOutPut(jobID string) string {
	if resp, err := client.Cli.Run("cat " + client.WorkDir + jobID + ".out"); err == nil {
		return string(resp)
	} else {
		return ""
	}
}

func (client *SSHClient) getError(jobID string) string {
	if resp, err := client.Cli.Run("cat " + client.WorkDir + jobID + ".err"); err == nil {
		return string(resp)
	} else {
		return ""
	}
}

/*--------------------SSH CONNECTION--------------------------*/
func newSSHClient(username, password string) *goph.Client {
	cli, err := goph.NewUnknown(username, gpuLoginAddr, goph.Password(password))
	if err != nil {
		fmt.Println("create SSH Client err: " + err.Error())
		return nil
	}
	return cli
}

func (client *SSHClient) Close() {
	err := client.Cli.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}

/*--------------------SHELL TOOLS--------------------------*/

func (client *SSHClient) Mkdir(dir string) (string, error) {
	cmd := fmt.Sprintf("mkdir %s", dir)
	resp, err := client.Cli.Run(cmd)
	return string(resp), err
}

func (client *SSHClient) CD(dir string) (string, error) {
	cmd := fmt.Sprintf("cd %s", dir)
	resp, err := client.Cli.Run(cmd)
	return string(resp), err
}

func (client *SSHClient) WriteFile(filename, content string) (string, error) {
	filePath := client.WorkDir + filename
	content = strings.Replace(content, "\"", "\\\"", -1)
	fmt.Println(content)
	cmd := fmt.Sprintf("echo \"%s\" > %s", content, filePath)
	resp, err := client.Cli.Run(cmd)
	return string(resp), err
}
