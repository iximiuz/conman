package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/kubelet/cri/streaming"

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
		"Root directory for persistent data, like container bundles, etc.")
	rootCmd.Flags().StringVarP(&cfg.RunRoot,
		"run-root", "n",
		config.DefaultRunRoot,
		"Root directory for runtime-only data, like sock & pid files.")
	rootCmd.Flags().StringVarP(&cfg.ContainerLogRoot,
		"container-logs", "L",
		config.DefaultContainerLogRoot,
		"Root directory for container logs.")
	rootCmd.Flags().StringVarP(&cfg.StreamingAddr,
		"streaming-addr", "S",
		config.DefaultStreamingAddr,
		"Network address (host:port) for streaming server (powers attach, exec, port-forwarding capabilities)")
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

		rs, err := cri.NewRuntimeService(
			oci.NewRuntime(
				fsutil.AssertExists(cfg.ShimmyPath),
				fsutil.AssertExists(cfg.RuntimePath),
				fsutil.EnsureExists(cfg.RuntimeRoot),
			),
			storage.NewContainerStore(fsutil.EnsureExists(cfg.LibRoot)),
			fsutil.EnsureExists(cfg.ContainerLogRoot),
			fsutil.EnsureExists(cfg.RunRoot, "exits"),
			fsutil.EnsureExists(cfg.RunRoot, "attach"),
		)
		if err != nil {
			logrus.Fatal(err)
		}

		sscfg := streaming.DefaultConfig
		sscfg.Addr = cfg.StreamingAddr
		ss, err := streaming.NewServer(sscfg, rs)
		if err != nil {
			logrus.Fatal(err)
		}
		go ss.Start(true)

		conman := server.New(rs, ss)
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
