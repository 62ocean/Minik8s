package kubectl

import (
	"github.com/urfave/cli/v2"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"log"
	"os"
)

var APIClient = HTTPClient.CreateHTTPClient(global.ServerHost)

func CmdExec() {
	app := &cli.App{
		Commands: []*cli.Command{
			CreateCmd(),
			DeleteCmd(),
			GetCmd(),
			DescribeCmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
