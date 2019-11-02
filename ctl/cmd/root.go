package cmd

import (
	"fmt"
	"os"

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
		fmt.Printf("Missed or unknown command.\n\n")
		cmd.Help()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
