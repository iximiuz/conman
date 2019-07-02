package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "conman",
	Short: "conman - simplistic container manager",
	Long:  `conman is a simplistic container manager, like CRI-O or containerd, but for edu purposes`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("conman root cmd")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
