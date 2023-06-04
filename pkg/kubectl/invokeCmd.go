package kubectl

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func InvokeCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "invoke",
		Usage: "invoke an object based on name",
		Subcommands: []*cli.Command{
			{
				Name:  "function",
				Usage: "invoke a function",
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return errors.New("the function name and the params must be specified")
					}
					name := c.Args().Get(0)
					params := c.Args().Get(1)
					fmt.Println("name: " + name)
					fmt.Println("params: " + params)
					paramjson, _ := json.Marshal(params)
					response := serverlessClient.Post("/invoke/function/"+name, paramjson)
					fmt.Println("result: " + response)
					return nil
				},
			},
			{
				Name:  "workflow",
				Usage: "invoke a workflow",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("the workflow name must be specified")
					}
					name := c.Args().Get(0)
					response := serverlessClient.Post("/invoke/workflow/"+name, nil)
					fmt.Println("result: " + response)
					return nil
				},
			},
		},
	}

	return cmd

}
