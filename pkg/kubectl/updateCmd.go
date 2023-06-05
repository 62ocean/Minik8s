package kubectl

import (
	"encoding/json"
	"errors"
	"github.com/urfave/cli/v2"
	"k8s/object"
	"k8s/pkg/util/parseYaml"
)

func UpdateCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "update",
		Usage: "update an object based on .yaml file",
		Subcommands: []*cli.Command{
			{
				Name:  "function",
				Usage: "update a function",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the file of a function",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the function name must be specified")
					}
					name := c.Args().First()
					filePath := c.String("f")

					var function object.Function
					function.Name = name
					function.Path = filePath
					funjson, _ := json.Marshal(function)
					serverlessClient.Post("/functions/update", funjson)
					return nil
				},
			},
			{
				Name:  "RS",
				Usage: "update a replicaset",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a replicaset",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
					rsJson, _ := json.Marshal(newRS)
					APIClient.Post("/replicasets/update", rsJson)
					return nil
				},
			},
			{
				Name:  "HPA",
				Usage: "update a HPA",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "f",
						Usage:    "the path of the configuration file of a HPA",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					filePath := c.String("f")
					newHPA := parseYaml.ParseYaml[object.Hpa](filePath)
					HPAJson, _ := json.Marshal(newHPA)
					APIClient.Post("/hpas/update", HPAJson)
					return nil
				},
			},
		},
	}

	return cmd

}
