package kubectl

import (
	"encoding/json"
	"errors"
	"github.com/urfave/cli/v2"
	"k8s/object"
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
		},
	}

	return cmd

}
