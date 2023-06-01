package kubectl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"k8s/object"
)

func DescribeCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "describe",
		Usage: "get the detailed running information of a pod or a service",
		Subcommands: []*cli.Command{
			{
				Name:  "pod",
				Usage: "get the running information of a pod",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the pod name must be specified")
					}
					podName := c.Args().First()
					podInfo := APIClient.Get("/pods/get/" + podName)
					podStorage := object.PodStorage{}
					_ = json.Unmarshal([]byte(podInfo), &podStorage)
					fmt.Println(podStorage)
					yamlData, err := yaml.Marshal(podStorage)
					if err != nil {
						fmt.Println("转换为 YAML 失败:", err)
						return nil
					}
					fmt.Println(string(yamlData))
					return nil
				},
			},
			{
				Name:  "service",
				Usage: "get the running information of a service",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the service name must be specified")
					}
					fmt.Println("service information: ", c.Args().First())
					serviceName := c.Args().First()

					getMsg, _ := json.Marshal(serviceName)
					resp := APIClient.Post("/services/get", getMsg)
					// fmt.Printf("Kubectl get service: %s\n", resp)
					service := object.Service{}
					var svStr string
					json.Unmarshal([]byte(resp), &svStr)
					json.Unmarshal([]byte(svStr), &service)

					yamlData, err := yaml.Marshal(service)
					if err != nil {
						fmt.Println("转换为 YAML 失败:", err)
						return nil
					}
					fmt.Println(string(yamlData))

					// fmt.Printf("Name: %s\n", service.Metadata.Name)
					// fmt.Printf("Type: %s\n", service.Spec.Type)
					// fmt.Printf("Selector:\n")
					// fmt.Printf("\tapp: %s\n", service.Spec.Selector.App)
					// fmt.Printf("\tenv: %s\n", service.Spec.Selector.Env)
					// fmt.Printf("ClusterIP: %s\n", service.Spec.ClusterIP)
					return nil
				},
			},
			{
				Name:  "dns",
				Usage: "get the information of dns and path",
				Action: func(c *cli.Context) error {
					resp := APIClient.Get("/dns/get")
					dns := object.Dns{}
					json.Unmarshal([]byte(resp), &dns)
					fmt.Println(dns)
					yamlData, err := yaml.Marshal(dns)
					if err != nil {
						fmt.Println("转换为 YAML 失败:", err)
						return nil
					}
					fmt.Println(string(yamlData))
					return nil
				},
			},
			{
				Name:  "replicaset",
				Usage: "get the detailed information of a replicaset",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the replicaset name must be specified")
					}
					name := c.Args().First()
					rsInfo := APIClient.Get("/replicasets/get/" + name)
					rs := object.ReplicaSet{}
					_ = json.Unmarshal([]byte(rsInfo), &rs)
					fmt.Println(rs)
					yamlData, err := yaml.Marshal(rs)
					if err != nil {
						fmt.Println("转换为 YAML 失败:", err)
						return nil
					}
					fmt.Println(string(yamlData))
					return nil
				},
			},
			{
				Name:  "hpa",
				Usage: "get the detailed information of a hpa",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the hpa name must be specified")
					}
					name := c.Args().First()
					hpaInfo := APIClient.Get("/hpas/get/" + name)
					hpa := object.Hpa{}
					_ = json.Unmarshal([]byte(hpaInfo), &hpa)
					fmt.Println(hpa)
					yamlData, err := yaml.Marshal(hpa)
					if err != nil {
						fmt.Println("转换为 YAML 失败:", err)
						return nil
					}
					fmt.Println(string(yamlData))
					return nil
				},
			},
		},
	}

	return cmd

}
