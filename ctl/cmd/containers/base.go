package containers

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iximiuz/conman/ctl/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(baseCmd)
}

type Options struct {
	Rootfs         string
	RootfsReadonly bool
	Command        string
}

var opts Options

var baseCmd = &cobra.Command{
	Use:   "container",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Missed or unknown container command.\n\n")
		cmd.Help()
	},
}
