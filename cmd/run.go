package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of conman",
	Long:  "Print the version number of conman",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("conman 0.0.1")
	},
}
