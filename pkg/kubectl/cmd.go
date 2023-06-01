package kubectl

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func CmdExec() {
	app := &cli.App{
		Commands: []*cli.Command{
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
				Name:  "create",
				Usage: "create a pod based on a pod.yaml",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a pod",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("create: ", c.String("f"))
					filePath := c.String("f")
					newPod := parseYaml.ParsePodYaml(filePath)
					// id, _ := uuid.NewUUID()
					// newPod.Metadata.Uid = id.String()
					client := HTTPClient.CreateHTTPClient(global.ServerHost)
					podJson, _ := json.Marshal(newPod)
					fmt.Println(newPod.Metadata.Name)
					client.Post("/pods/create", podJson)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete a pod based on the type and name of a pod.yaml",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a pod",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("delete: ", c.String("f"))
					//apiserver.CreatePod()
					return nil
				},
			},
			{
				Name:  "describe",
				Usage: "get the running information of a pod or a service",
				Subcommands: []*cli.Command{
					{
						Name:  "pod",
						Usage: "get the running information of a pod",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the pod name must be specified")
							}
							fmt.Println("pod information: ", c.Args().First())
							//apiserver.DescribePod()
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
							//apiserver.DescribeService()
							return nil
						},
					},
				},
			},
			//{
			//	Name:    "template",
			//	Aliases: []string{"t"},
			//	Usage:   "options for task templates",
			//	Subcommands: []*cli.Command{
			//		{
			//			Name:  "add",
			//			Usage: "add a new template",
			//			Action: func(c *cli.Context) error {
			//				fmt.Println("new task template: ", c.Args().First())
			//				return nil
			//			},
			//		},
			//		{
			//			Name:  "remove",
			//			Usage: "remove an existing template",
			//			Action: func(c *cli.Context) error {
			//				fmt.Println("removed task template: ", c.Args().First())
			//				return nil
			//			},
			//		},
			//	},
			//},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
