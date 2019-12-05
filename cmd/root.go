package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/cri"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/server"
)

var cfg config.Config

func init() {
	rootCmd.Flags().StringVarP(&cfg.Listen,
		"listen", "l",
		config.DefaultListen,
		"Daemon listen address")
	rootCmd.Flags().StringVarP(&cfg.LibRoot,
		"lib-root", "b",
		config.DefaultLibRoot,
		"TODO: ...")
	rootCmd.Flags().StringVarP(&cfg.RunRoot,
		"run-root", "n",
		config.DefaultRunRoot,
		"TODO: ...")
	rootCmd.Flags().StringVarP(&cfg.ShimmyPath,
		"shimmy-path", "s",
		config.DefaultShimmyPath,
		"Path to OCI runtime shime executable (shimmy)")
	rootCmd.Flags().StringVarP(&cfg.RuntimePath,
		"runtime-path", "r",
		config.DefaultRuntimePath,
		"Path to OCI-compatible runtime executable")
	rootCmd.Flags().StringVarP(&cfg.RuntimeRoot,
		"runtime-root", "t",
		config.DefaultRuntimeRoot,
		"OCI runtime root directory")

	// TODO: configure it
	logrus.SetLevel(logrus.TraceLevel)
}

var rootCmd = &cobra.Command{
	Use:   "conmand",
	Short: "conmand - simplistic container manager",
	Long: `conmand is a simplistic container manager, 
like CRI-O or containerd, but for edu purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Conman's here!")

		assertExists(cfg.RuntimePath)
		ensureExists(cfg.LibRoot)

		rs, err := cri.NewRuntimeService(
			oci.NewRuntime(
				cfg.ShimmyPath,
				cfg.RuntimePath,
				cfg.RuntimeRoot,
			),
			storage.NewContainerStore(cfg.LibRoot),
		)
		if err != nil {
			logrus.Fatal(err)
		}

		conman := server.New(rs)
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

func assertExists(filename string) {
	ok, err := fsutil.Exists(filename)
	if !ok || err != nil {
		logrus.WithError(err).Fatal("File is not reachable: " + filename)
	}
}

func ensureExists(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		logrus.WithError(err).Fatal("File is not reachable: " + dir)
	}
}
