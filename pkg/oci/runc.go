package oci

import (
	"os"
	"os/exec"

	"github.com/iximiuz/conman/pkg/container"
)

type runcRuntime struct {
	exePath  string
	rootPath string
}

func NewRuntime(exePath, rootPath string) Runtime {
	return &runcRuntime{
		exePath:  exePath,
		rootPath: rootPath,
	}
}

func (r *runcRuntime) CreateContainer(id container.ID, bundle string) error {
	cmd := exec.Command(
		r.exePath,
		"-bundle", bundle,
		string(id),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func (r *runcRuntime) StartContainer() {
	panic("not implemented")
}

func (r *runcRuntime) KillContainer() {
	panic("not implemented")
}

func (r *runcRuntime) DeleteContainer() {
	panic("not implemented")
}

func (r *runcRuntime) ContainerState() {
	panic("not implemented")
}
