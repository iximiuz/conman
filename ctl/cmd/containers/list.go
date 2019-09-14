package containers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	cmdutil "github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	baseCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connect()
		defer conn.Close()

		resp, err := client.ListContainers(
			context.Background(),
			&server.ListContainersRequest{},
		)
		if err != nil {
			logrus.WithError(err).
				Fatal("Command failed (see conmand logs for details)")
		}
		cmdutil.Print(resp)
	},
}
