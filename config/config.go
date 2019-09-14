package config

const (
	DefaultLibRoot     = "/var/lib/conman"
	DefaultListen      = "/run/conmand.sock"
	DefaultRunRoot     = "/run/conman"
	DefaultRuntimePath = "/usr/bin/runc"
	DefaultRuntimeRoot = "/run/conman-runc"
)

type Config struct {
	Listen string

	// Root directory to store long living data (images, containers, etc).
	LibRoot string

	// Root directory to store state of the conman daemon.
	RunRoot string

	RuntimePath string

	RuntimeRoot string
}

func TestConfig() (*Config, error) {
	cfg := &Config{
		LibRoot:     DefaultLibRoot,
		Listen:      DefaultListen,
		RunRoot:     DefaultRunRoot,
		RuntimePath: DefaultRuntimePath,
		RuntimeRoot: DefaultRuntimeRoot,
	}
	return cfg, nil
}
