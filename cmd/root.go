package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/cri"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/server"
)

var cfg config.Config

func init() {
	rootCmd.Flags().StringVarP(&cfg.Listen,
		"listen", "l",
		"/run/conmand.sock",
		"Daemon listen address")
	rootCmd.Flags().StringVarP(&cfg.LibRoot,
		"lib-root", "b",
		"/var/lib/conman",
		"TODO: ...")
	rootCmd.Flags().StringVarP(&cfg.RunRoot,
		"run-root", "n",
		"/run/conman",
		"TODO: ...")
	rootCmd.Flags().StringVarP(&cfg.RuntimePath,
		"runtime-path", "r",
		"/usr/bin/runc",
		"Path to OCI-compatible runtime executable")
	rootCmd.Flags().StringVarP(&cfg.RuntimeRoot,
		"runtime-root", "t",
		"/run/runc",
		"OCI runtime root directory")
}

var rootCmd = &cobra.Command{
	Use:   "conman",
	Short: "conman - simplistic container manager",
	Long: `conman is a simplistic container manager, 
like CRI-O or containerd, but for edu purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Conman's here!")

		conman := server.New(
			cri.NewRuntimeService(
				oci.NewRuntime(
					cfg.RuntimePath,
					cfg.RuntimeRoot,
				),
				storage.New(cfg.LibRoot),
			),
		)
		if err := conman.Serve("unix", cfg.Listen); err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Conman has left the building!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
