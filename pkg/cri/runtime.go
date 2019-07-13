package cri

import (
	"github.com/sirupsen/logrus"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/rollback"
	"github.com/iximiuz/conman/pkg/storage"
)

// RuntimeService is a CRI-alike service
// to manage container & sandbox runtimes.
type RuntimeService interface {
	CreateContainer(ContainerOptions) (*container.Container, error)
	StartContainer(container.ID) error
	StopContainer(container.ID) error
	// ExecSync
	// Attach

	// RunPodSandbox
	// StopPodSandbox
	// RemovePodSandbox
	// PodSandboxStatus
	// ListPodSandbox
}

type runtimeService struct {
	runtime oci.Runtime
	cmap    *container.Map
	cstore  storage.ContainerStore
}

func NewRuntimeService(
	runtime oci.Runtime,
	cstore storage.ContainerStore,
) RuntimeService {
	return &runtimeService{
		runtime: runtime,
		cmap:    container.NewMap(),
		cstore:  cstore,
	}
}

func (s *runtimeService) CreateContainer(
	opts ContainerOptions,
) (cont *container.Container, err error) {
	rb := rollback.New()
	defer func() { _ = err == nil || rb.Execute() }()

	cont, err = container.New(
		container.RandID(),
		opts.Name,
	)
	if err != nil {
		return
	}

	if err = s.cmap.Add(cont, rb); err != nil {
		return
	}

	if err = s.cstore.CreateContainer(cont, rb); err != nil {
		return
	}

	// TODO: spec := oci.NewSpec(...)
	var spec interface{} = nil
	err = s.cstore.CreateContainerBundle(cont.ID(), spec, opts.RootfsPath)
	if err != nil {
		return
	}

	// TODO: finally launch runc
	return
}

func (r *runtimeService) StartContainer(
	_id container.ID,
) error {
	logrus.Debug("StartContainer")
	return nil
}

func (r *runtimeService) StopContainer(
	_id container.ID,
) error {
	logrus.Debug("StopContainer")
	return nil
}

type ContainerOptions struct {
	Name           string
	RootfsPath     string
	RootfsReadonly bool
}
