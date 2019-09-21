package containers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	cmdutil "github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	baseCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start <container-id>",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := cmdutil.Connect()
		defer conn.Close()

		resp, err := client.StartContainer(
			context.Background(),
			&server.StartContainerRequest{
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
