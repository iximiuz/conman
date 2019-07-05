package config

type Config struct {
	Listen string

	// Root directory to store long living data (images, containers, etc).
	LibRoot string

	// Root directory to store state of the conman daemon.
	RunRoot string

	RuntimePath string

	RuntimeRoot string
}
