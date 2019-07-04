package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/iximiuz/conman/runtime"
	"github.com/iximiuz/conman/server"
)

var runtimePath string
var listen string

func init() {
	rootCmd.Flags().StringVarP(&listen, "listen", "l", "",
		"Daemon listen address")
	rootCmd.Flags().StringVarP(&runtimePath, "runtime", "r", "",
		"Path to OCI-compatible runtime executable")
}

var rootCmd = &cobra.Command{
	Use:   "conman",
	Short: "conman - simplistic container manager",
	Long: `conman is a simplistic container manager, 
like CRI-O or containerd, but for edu purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting conman...")

		conman := server.New(runtime.NewRunc(runtimePath))
		if err := conman.Serve("unix", listen); err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Conman has left the building!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
