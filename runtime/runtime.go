package runtime

type Runtime interface {
	CreateContainer(name string) (*Container, error)
	StartContainer(id ContainerID) error
}
