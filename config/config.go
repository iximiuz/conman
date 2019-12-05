package config

import (
	"flag"
)

const (
	DefaultLibRoot     = "/var/lib/conman"
	DefaultListen      = "/run/conmand.sock"
	DefaultRunRoot     = "/run/conman"
	DefaultShimmyPath  = "/usr/local/bin/shimmy"
	DefaultRuntimePath = "/usr/bin/runc"
	DefaultRuntimeRoot = "/run/conman-runc"
)

type Config struct {
	Listen string

	// Root directory to store long living data (images, containers, etc).
	LibRoot string

	// Root directory to store state of the conman daemon.
	RunRoot string

	// Path to OCI runtime shim executable, aka shimmy.
	ShimmyPath string

	RuntimePath string

	RuntimeRoot string
}

func TestConfigFromFlags() *Config {
	cfg := &Config{
		LibRoot:     DefaultLibRoot,
		Listen:      DefaultListen,
		RunRoot:     DefaultRunRoot,
		RuntimeRoot: DefaultRuntimeRoot,
	}

	flag.StringVar(
		&cfg.RuntimePath,
		"runtime",
		DefaultRuntimePath,
		"Path to runc executable file",
	)
	flag.StringVar(
		&cfg.ShimmyPath,
		"shimmy",
		DefaultShimmyPath,
		"Path to shimmy executable file",
	)
	flag.Parse()

	return cfg
}
