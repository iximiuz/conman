package containers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	cmd.RootCmd.AddCommand(baseCmd)
}

type Options struct {
	Rootfs         string
	RootfsReadonly bool
	Command        string
}

var opts Options

var baseCmd = &cobra.Command{
	Use:   "container",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Fatal("container action required")
	},
}

func connect() (server.ConmanClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("unix://"+cmd.OptHost, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	return server.NewConmanClient(conn), conn
}
