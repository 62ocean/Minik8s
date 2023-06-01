package kubectl

import (
	"encoding/json"
	"errors"
	"github.com/urfave/cli/v2"
	"k8s/object"
	"k8s/pkg/util/parseYaml"
	"log"
)

func DeleteCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "delete",
		Usage: "delete an object based on name",
		Subcommands: []*cli.Command{
			{
				Name:  "pod",
				Usage: "delete a pod",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the pod name must be specified")
					}
					name := c.Args().First()
					nameReq, _ := json.Marshal(name)
					APIClient.Post("/pods/remove", nameReq)
					return nil
				},
			},
			{
				Name:  "service",
				Usage: "delete a service",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the service name must be specified")
					}
					name := c.Args().First()
					nameReq, _ := json.Marshal(name)
					APIClient.Post("/services/remove", nameReq)
					return nil
				},
			},
			{
				Name:  "RS",
				Usage: "get the running information of a replicaset",
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					log.Println("delete RS: ", c.String("f"))
					newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
					rsJson, _ := json.Marshal(newRS)
					log.Println(newRS)
					APIClient.Post("/replicasets/delete", rsJson)
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
					log.Println("delete HPA: ", c.String("f"))
					newHPA := parseYaml.ParseYaml[object.Hpa](filePath)
					HPAJson, _ := json.Marshal(newHPA)
					log.Println(newHPA)
					APIClient.Post("/hpas/delete", HPAJson)
					return nil
				},
			},
		},
	}

	return cmd

}
