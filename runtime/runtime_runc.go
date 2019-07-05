package runtime

import (
	"github.com/davecgh/go-spew/spew"
	specgenerate "github.com/opencontainers/runtime-tools/generate"
	"github.com/satori/go.uuid"
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
	id := reserveContID()
	cont, err := NewContainer(id, name)
	// TODO: save ID & name, check for uniqueness
	if err != nil {
		return nil, err
	}
	// TODO: defer clean up on error

	// TODO: create container dir in /run/conman
	// TODO: copy rootfs to container dir

	gen, err := specgenerate.New("linux")
	if err != nil {
		return nil, err
	}
	gen.HostSpecific = true
	// TODO: set cmd & args from container params
	// TODO: save config.json to container dir
	logrus.Info(spew.Sdump(gen))

	// TODO: finally launch runc

	return cont, err
}

func (r *runcRuntime) StartContainer(
	id ContainerID,
) error {
	logrus.Debug("StartContainer")
	return nil
}

func reserveContID() ContainerID {
	return ContainerID(uuid.NewV4().String())
}
