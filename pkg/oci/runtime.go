package oci

import "github.com/iximiuz/conman/pkg/container"

// Runtime represents an OCI container runtime interface.
type Runtime interface {
	CreateContainer(id container.ID, bundle string) error
	StartContainer(id container.ID) error
	KillContainer()
	DeleteContainer()
	ContainerState()
}
