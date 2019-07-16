package oci

import "github.com/iximiuz/conman/pkg/container"

type Runtime interface {
	CreateContainer(id container.ID, bundle string) error
	StartContainer()
	KillContainer()
	DeleteContainer()
	ContainerState()
}
