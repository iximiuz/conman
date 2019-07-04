package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/iximiuz/conman/server"
)

func init() {
	rootCmd.AddCommand(createContainerCmd)
}

var createContainerCmd = &cobra.Command{
	Use:   "create-container",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("unix://"+optHost, grpc.WithInsecure())
		if err != nil {
			logrus.Fatal(err)
		}
		defer conn.Close()

		client := server.NewConmanClient(conn)
		resp, err := client.CreateContainer(
			context.Background(),
			&server.CreateContainerRequest{},
		)
		logrus.Info(resp, err)
	},
}
