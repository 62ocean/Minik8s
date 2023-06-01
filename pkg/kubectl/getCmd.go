package kubectl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"k8s/object"
	"time"
)

func GetCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "get",
		Usage: "get the running information of a pod or a service",
		Subcommands: []*cli.Command{
			{
				Name:  "pod",
				Usage: "get the running information of a pod",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the pod name must be specified")
					}
					name := c.Args().First()
					podInfo := APIClient.Get("/pods/get/" + name)
					podStorage := object.PodStorage{}
					_ = json.Unmarshal([]byte(podInfo), &podStorage)
					fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
					createTime := podStorage.Config.Metadata.CreationTimestamp
					newtime := time.Now()
					d := newtime.Sub(createTime)
					fmt.Printf("%s\t\t\t%s\t\t\t%s\n", name, podStorage.Status.ToString(), d.Truncate(time.Second).String())
					return nil
				},
			},
			{
				Name:  "pods",
				Usage: "get the running information of all pod",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/pods/getAll")
					var podList map[string]string
					_ = json.Unmarshal([]byte(response), &podList)
					fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
					for _, val := range podList {
						podStorage := object.PodStorage{}
						_ = json.Unmarshal([]byte(val), &podStorage)
						createTime := podStorage.Config.Metadata.CreationTimestamp
						newtime := time.Now()
						d := newtime.Sub(createTime)
						fmt.Printf("%s\t\t\t%s\t\t\t%s\n", podStorage.Config.Metadata.Name, podStorage.Status.ToString(), d.Truncate(time.Second).String())

					}
					return nil
				},
			},
			{
				Name:  "GPUJob",
				Usage: "get the running information of a GPUJob",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the job name must be specified")
					}
					name := c.Args().First()
					jobInfo := APIClient.Get("/gpuJobs/get/" + name)
					job := object.GPUJob{}
					_ = json.Unmarshal([]byte(jobInfo), &job)
					fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
					createTime := job.Metadata.CreationTimestamp
					newtime := time.Now()
					d := newtime.Sub(createTime)
					fmt.Printf("%s\t\t\t%s\t\t\t%s\n", name, job.Status.ToString(), d.Truncate(time.Second).String())
					if job.Status == 3 {
						fmt.Printf("OUTPUT: \n")
						fmt.Println(job.Output)
					}
					return nil
				},
			},
			{
				Name:  "services",
				Usage: "get the running information of all services",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/services/getAll")
					var serviceList map[string]string
					_ = json.Unmarshal([]byte(response), &serviceList)
					fmt.Println("NAME\t\t\tCLUSTERIP\t\t\tLABEL")
					for _, val := range serviceList {
						service := object.Service{}
						_ = json.Unmarshal([]byte(val), &service)
						label := fmt.Sprint("app:%s env:%s", service.Metadata.Labels.App, service.Metadata.Labels.Env)
						fmt.Printf("%s\t\t\t%s\t\t\t%s\n", service.Metadata.Name, service.Spec.ClusterIP, label)
					}
					return nil
				},
			},
			{
				Name:  "replicasets",
				Usage: "list all replicasets",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/replicasets/getAll")
					var replicasetList map[string]string
					_ = json.Unmarshal([]byte(response), &replicasetList)
					fmt.Println("NAME\t\t\tREPLICA\t\t\tLABEL")
					for _, val := range replicasetList {
						rs := object.ReplicaSet{}
						_ = json.Unmarshal([]byte(val), &rs)
						label := fmt.Sprint("app:%s env:%s", rs.Spec.Selector.MatchLabels.App, rs.Spec.Selector.MatchLabels.Env)
						fmt.Printf("%s\t\t\t%s\t\t\t%s\n", rs.Metadata.Name, rs.Spec.Replicas, label)
					}
					return nil
				},
			},
			{
				Name:  "hpas",
				Usage: "list all hpas",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/hpas/getAll")
					var hpaList map[string]string
					_ = json.Unmarshal([]byte(response), &hpaList)
					fmt.Println("NAME\t\t\tMIN REPLICA\t\t\tMAX REPLICA\t\t\tCPU METRIC\t\t\tMEMORY METRIC")
					for _, val := range hpaList {
						hpa := object.Hpa{}
						_ = json.Unmarshal([]byte(val), &hpa)
						//label := fmt.Sprint("app:%s env:%s", rs.Spec.Selector.MatchLabels.App, rs.Spec.Selector.MatchLabels.Env)
						fmt.Printf("%s\t\t\t%s\t\t\t%s\t\t\t%s\t\t\t%s\n", hpa.Metadata.Name, hpa.Spec.MinReplicas, hpa.Spec.MaxReplicas,
							hpa.Spec.Metrics[0].Resource.Target.AverageUtilization, hpa.Spec.Metrics[1].Resource.Target.AverageUtilization)
					}
					return nil
				},
			},
			{
				Name:  "functions",
				Usage: "list all functions",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/functions/getAll")
					var functionList map[string]string
					_ = json.Unmarshal([]byte(response), &functionList)
					fmt.Println("NAME\t\t\tPATH\t\t\tIMAGE")
					for _, val := range functionList {
						fun := object.Function{}
						_ = json.Unmarshal([]byte(val), &fun)
						//label := fmt.Sprint("app:%s env:%s", rs.Spec.Selector.MatchLabels.App, rs.Spec.Selector.MatchLabels.Env)
						fmt.Printf("%s\t\t\t%s\t\t\t%s\n", fun.Name, fun.Path, fun.Image)
					}
					return nil
				},
			},
			{
				Name:  "workflows",
				Usage: "list all workflows",
				Action: func(c *cli.Context) error {
					response := APIClient.Get("/workflows/getAll")
					var wfList map[string]string
					_ = json.Unmarshal([]byte(response), &wfList)
					fmt.Println("NAME")
					for _, val := range wfList {
						wf := object.Workflow{}
						_ = json.Unmarshal([]byte(val), &wf)
						//label := fmt.Sprint("app:%s env:%s", rs.Spec.Selector.MatchLabels.App, rs.Spec.Selector.MatchLabels.Env)
						fmt.Printf("%s\n", wf.Metadata.Name)
					}
					return nil
				},
			},
		},
	}

	return cmd

}
