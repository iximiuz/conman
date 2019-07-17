package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var OptHost string

func init() {
	RootCmd.PersistentFlags().StringVarP(&OptHost,
		"host", "H",
		"/run/conmand.sock",
		"Daemon socket to connect")
}

var RootCmd = &cobra.Command{
	Use:   "conmanctl",
	Short: "conmanctl - CLI tool to communicate with conmand",
	Long:  `conmanctl - CLI tool to communicate with conmand.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Fatal("action required")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
