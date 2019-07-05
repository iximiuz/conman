package runtime

import (
	"github.com/sirupsen/logrus"
)

type runcRuntime struct {
	exePath string
}

func NewRunc(exePath string) Runtime {
	return &runcRuntime{
		exePath: exePath,
	}
}

func (r *runcRuntime) CreateContainer(
	name string,
) (*Container, error) {
	id := ContainerID("cont_foobar_123")
	cont, err := NewContainer(id, name)
	return cont, err
}

func (r *runcRuntime) StartContainer(
	id ContainerID,
) error {
	logrus.Debug("StartContainer")
	return nil
}
