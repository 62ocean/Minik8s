package kubectl

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/urfave/cli/v2"
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
					log.Println("delete pod " + name)
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
				Usage: "delete a replicaset",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the replicaset name must be specified")
					}
					name := c.Args().First()
					APIClient.Del("/replicasets/remove/" + name)
					return nil
				},
			},
			{
				Name:  "HPA",
				Usage: "delete a hpa",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the hpa name must be specified")
					}
					name := c.Args().First()
					APIClient.Del("/hpas/remove/" + name)
					return nil
				},
			},
			{
				Name:  "function",
				Usage: "delete a function",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the function name must be specified")
					}
					name := c.Args().First()
					serverlessClient.Del("/functions/remove/" + name)
					return nil
				},
			},
		},
	}

	return cmd

}
