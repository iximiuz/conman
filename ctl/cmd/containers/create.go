package containers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	cmdutil "github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	createCmd.PersistentFlags().StringVarP(&opts.Rootfs,
		"image", "i",
		"",
		"Container rootfs image (required)")
	createCmd.MarkPersistentFlagRequired("image")

	createCmd.PersistentFlags().BoolVarP(&opts.RootfsReadonly,
		"rootfs-readonly", "R",
		true,
		"Wether container can modify its rootfs")

	baseCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [command options] <container-name> -- <command> [args...]",
	Short: "",
	Long:  "",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := cmdutil.Connect()
		defer conn.Close()

		resp, err := client.CreateContainer(
			context.Background(),
			&server.CreateContainerRequest{
				Name:           args[0],
				RootfsPath:     opts.Rootfs,
				RootfsReadonly: opts.RootfsReadonly,
				Command:        args[1],
				Args:           args[2:],
			},
		)
		if err != nil {
			logrus.WithError(err).
				Fatal("Command failed (see conmand logs for details)")
		}
		cmdutil.Print(resp)
	},
}
