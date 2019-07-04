package runtime

type Runtime interface {
	CreateContainer(bundlePath string) (*Container, error)
	StartContainer(bundlePath string) (*Container, error)
}
