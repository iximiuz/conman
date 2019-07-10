package oci

type Runtime interface {
	CreateContainer()
	StartContainer()
	KillContainer()
	DeleteContainer()
	ContainerState()
}
