package oci

import "github.com/iximiuz/conman/pkg/container"

type Runtime interface {
	CreateContainer(id container.ID, bundle string) error
	StartContainer(id container.ID) error
	KillContainer()
	DeleteContainer()
	ContainerState()
}
