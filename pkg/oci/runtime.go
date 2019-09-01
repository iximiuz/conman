package oci

import (
	"os"

	"github.com/iximiuz/conman/pkg/container"
)

// Runtime represents an OCI container runtime interface.
type Runtime interface {
	CreateContainer(id container.ID, bundle string) error
	StartContainer(id container.ID, bundle string) error
	KillContainer(id container.ID, sig os.Signal) error
	DeleteContainer()
	ContainerState(container.ID) (StateResp, error)
}

type StateResp struct {
	Pid    int    `json:"pid"`
	Status string `json:"status"`
}
