package GPU

import (
	"fmt"
	"github.com/melbahja/goph"
	"os/exec"
	"runtime"
)

const (
	gpuUser      = "stu1639"
	gpuPasswd    = "34j3oQ*%"
	gpuLoginAddr = "pilogin.hpc.sjtu.edu.cn"
	gpuDataAddr  = "data.hpc.sjtu.edu.cn"
	accountType  = "acct-stu"
)

type SSHClient struct {
	Username string
	Password string
	Cli      *goph.Client
}

func newSSHClient(username, password string) *goph.Client {
	cli, err := goph.NewUnknown(username, gpuLoginAddr, goph.Password(password))
	if err != nil {
		fmt.Println("create SSH Client err: " + err.Error())
		return nil
	}
	return cli
}

func NewClient() *SSHClient {
	sshCli := newSSHClient(gpuUser, gpuPasswd)
	return &SSHClient{
		Username: gpuUser,
		Password: gpuPasswd,
		Cli:      sshCli,
	}
}

func (client *SSHClient) Scp(sourcePath string, dstPath string) error {
	if runtime.GOOS == "linux" {
		remoteAddr := fmt.Sprintf("%s@%s:%s", client.Username, gpuDataAddr, dstPath)
		cmd := exec.Command("scp", "-r", sourcePath, remoteAddr)
		return cmd.Run()
	}
	return fmt.Errorf("scp is not supported in your os")
}
