package kubectl

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func CmdExec() {
	app := &cli.App{
		Commands: []*cli.Command{
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
					//apiserver.CreatePod()
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
