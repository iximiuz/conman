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
		"image", "I",
		"",
		"Container rootfs image (required)")
	createCmd.MarkPersistentFlagRequired("image")

	createCmd.PersistentFlags().BoolVarP(&opts.RootfsReadonly,
		"rootfs-readonly", "R",
		true,
		"Wether container can modify its rootfs")

	createCmd.PersistentFlags().BoolVarP(&opts.Stdin,
		"stdin", "i",
		false,
		"Keep container's STDIN open (interactive mode)")

	createCmd.PersistentFlags().BoolVarP(&opts.LeaveStdinOpen,
		"leave-stdin-open", "",
		false,
		"Leave container's STDIN open after first attach session completes")

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
				Stdin:          opts.Stdin,
				StdinOnce:      !opts.LeaveStdinOpen,
			},
		)
		if err != nil {
			logrus.WithError(err).
				Fatal("Command failed (see conmand logs for details)")
		}
		cmdutil.Print(resp)
	},
}
