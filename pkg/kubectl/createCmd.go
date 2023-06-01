package kubectl

import (
	"encoding/json"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func RunCommand(cmd string) {
	fmt.Printf("RunCmd: %s\n", cmd)
	command := exec.Command("/bin/bash", "-c", cmd)
	// out, err := command.Output()
	// if err != nil {
	// 	fmt.Println("output : ")
	// 	fmt.Println(out)
	// }
	out, err := command.CombinedOutput()
	log.Printf("out: %s", string(out))
	if err != nil {
		fmt.Printf("ERROR: run cmd error: %s\n", err.Error())
		// panic("ERROR: " + err.Error())
	}
	// return string(out)
}
func CreateCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "create",
		Usage: "create an object based on .yaml file",
		Subcommands: []*cli.Command{
			{
				Name:  "node",
				Usage: "create a node",
				Action: func(c *cli.Context) error {
					fmt.Printf("prepare for environment...\n")
					RunCommand("make clean-env")
					RunCommand("make kill-all")
					fmt.Printf("build code...\n")
					RunCommand("make node")
					fmt.Printf("start node...\n")
					RunCommand("make node_start")
					return nil
				},
			},

			{
				Name:  "dns",
				Usage: "create a dns based on a dns.yaml",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a dns",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("create: ", c.String("f"))
					filePath := c.String("f")
					newPod := parseYaml.ParseYaml[object.Dns](filePath)
					// id, _ := uuid.NewUUID()
					// newPod.Metadata.Uid = id.String()
					client := HTTPClient.CreateHTTPClient(global.ServerHost)
					dnsJson, _ := json.Marshal(newPod)
					fmt.Println(newPod.Metadata.Name)
					client.Post("/dns/create", dnsJson)
					return nil
				},
			},
			{
				Name:  "pod",
				Usage: "create a pod",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a pod",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("create pod: ", c.String("f"))
					newPod := parseYaml.ParseYaml[object.Pod](filePath)
					podJson, _ := json.Marshal(newPod)
					log.Println(newPod)
					APIClient.Post("/pods/create", podJson)
					return nil
				},
			},
			{
				Name:  "service",
				Usage: "create a service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a service",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("create service: ", c.String("f"))
					newService := parseYaml.ParseYaml[object.Service](filePath)
					serviceJson, _ := json.Marshal(newService)
					log.Println(newService)
					APIClient.Post("/services/create", serviceJson)
					return nil
				},
			},
			{
				Name:  "RS",
				Usage: "get the running information of a replicaset",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a replicaset",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("create RS: ", c.String("f"))
					newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
					rsJson, _ := json.Marshal(newRS)
					log.Println(newRS)
					APIClient.Post("/replicasets/create", rsJson)
					return nil
				},
			},
			{
				Name:  "HPA",
				Usage: "get the running information of a HPA",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a HPA",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("create HPA: ", c.String("f"))
					newHPA := parseYaml.ParseYaml[object.Hpa](filePath)
					HPAJson, _ := json.Marshal(newHPA)
					log.Println(newHPA)
					APIClient.Post("/hpas/create", HPAJson)
					return nil
				},
			},
			{
				Name:  "GPUJob",
				Usage: "get the running information of a service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a pod",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("create GPUJob: ", c.String("f"))
					// job存入apiserver
					job := parseYaml.ParseYaml[object.GPUJob](filePath)
					job.Status = object.PENDING
					jobInfo, _ := json.Marshal(job)
					APIClient.Post("/gpuJobs/create", jobInfo)

					// 构造pod 存入apiserver
					port := object.ContainerPort{Port: 8080}
					container := object.Container{
						Name:  "commit_" + "GPUJob_" + job.Metadata.Name,
						Image: "saltfishy/gpu_server:v8",
						Ports: []object.ContainerPort{
							port,
						},
						Command: []string{
							"./main ",
						},
						Args: []string{
							job.Metadata.Name,
						},
						// TODO 此处写入kubectl时需要修改为参数指定的文件路径
						CopyFile: job.Spec.Program,
						CopyDst:  "/apps",
					}
					newPod := object.Pod{
						ApiVersion: "v1",
						Kind:       "Pod",
						Metadata: object.Metadata{
							Name: "GPUJob_" + job.Metadata.Name,
							Labels: object.Labels{
								App: "GPU",
								Env: "prod",
							},
						},
						Spec: object.PodSpec{
							Containers: []object.Container{
								container,
							},
						},
					}
					podInfo, _ := json.Marshal(newPod)
					APIClient.Post("/pods/create", podInfo)
					return nil
				},
			},
		},
	}

	return cmd

}
