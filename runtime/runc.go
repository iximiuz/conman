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
	_bundlePath string,
) (*Container, error) {
	logrus.Debug("CreateContainer")
	return &Container{
		id: "create_cont_foobar",
	}, nil
}

func (r *runcRuntime) StartContainer(
	bundlePath string,
) (*Container, error) {
	logrus.Debug("StartContainer")
	return &Container{
		id: "start_cont_foobar",
	}, nil
}
