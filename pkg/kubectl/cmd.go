package kubectl

import (
	"github.com/urfave/cli/v2"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"log"
	"os"
)

var APIClient = HTTPClient.CreateHTTPClient(global.ServerHost)
var serverlessClient = HTTPClient.CreateHTTPClient(global.ServerlessHost)

func CmdExec() {
	app := &cli.App{
		Commands: []*cli.Command{
			CreateCmd(),
			DeleteCmd(),
			GetCmd(),
			DescribeCmd(),
			UpdateCmd(),
			InvokeCmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
