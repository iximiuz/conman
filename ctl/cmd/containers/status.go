package containers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	cmdutil "github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	baseCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status <container-id>",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connect()
		defer conn.Close()

		resp, err := client.ContainerStatus(
			context.Background(),
			&server.ContainerStatusRequest{
				ContainerId: args[0],
			},
		)
		if err != nil {
			logrus.WithError(err).
				Fatal("Command failed (see conmand logs for details)")
		}
		cmdutil.Print(resp)
	},
}
