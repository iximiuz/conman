package containers

import (
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	cmdutil "github.com/iximiuz/conman/ctl/cmd"
	"github.com/iximiuz/conman/server"
)

func init() {
	baseCmd.AddCommand(attachCmd)
}

var attachCmd = &cobra.Command{
	Use:   "attach <container-id>",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := cmdutil.Connect()
		defer conn.Close()

		resp, err := client.Attach(
			context.Background(),
			&server.AttachRequest{
				ContainerId: args[0],
				Tty:         false,
				Stdin:       true,
				Stdout:      true,
				Stderr:      true,
			},
		)
		if err != nil {
			logrus.WithError(err).
				Fatal("Command failed (see conmand logs for details)")
		}

		url, err := url.Parse(resp.Url)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to parse stream URL")
		}

		executor, err := remotecommand.NewSPDYExecutor(
			&rest.Config{
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			},
			"POST",
			url,
		)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create stream executor")
		}

		streamOptions := remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Tty:    false,
		}
		if err := executor.Stream(streamOptions); err != nil {
			logrus.WithError(err).Fatal("executor.Stream() failed")
		}
	},
}
