package cri

import (
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/rollback"
	"github.com/iximiuz/conman/pkg/storage"
)

// RuntimeService is a service to manage container & sandbox runtimes.
// While it resembles the CRI runtime interface, it does not follow it
// strictly. The purpose of this service is to support the public-facing
// CRI runtime service (see server.Server).
type RuntimeService interface {
	// CreateContainer prepares a new container bundle on disk
	// and starts runc init, but does not start a specified process.
	CreateContainer(ContainerOptions) (*container.Container, error)

	// StartContainer actually starts a pre-defined process in
	// a container created via CreateContainer() call.
	StartContainer(container.ID) error

	// StopContainer signals the container to finish itself.
	StopContainer(id container.ID, timeout time.Duration) error

	// +RemoveContainer(container.ID) error

	// +ListContainers()

	// GetContainer returns the container doing a state request
	// from the OCI runtime if applicable.
	GetContainer(container.ID) (*container.Container, error)

	// UpdateContainerResources
	// ReopenContainerLog
	// ExecSync
	// Exec
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
	// TODO: add mutex lock to the method

	rb := rollback.New()
	defer func() { _ = err == nil || rb.Execute() }()

	cont, err = container.New(
		container.RandID(),
		opts.Name,
	)
	if err != nil {
		return
	}
	cont.SetStatus(container.StatusCreated)

	if err = s.cmap.Add(cont, rb); err != nil {
		return
	}

	h, err := s.cstore.CreateContainer(cont, rb)
	if err != nil {
		return
	}

	spec, err := oci.NewSpec(oci.SpecOptions{
		Command:      opts.Command,
		Args:         opts.Args,
		RootPath:     h.RootfsDir(),
		RootReadonly: opts.RootfsReadonly,
	})
	if err != nil {
		return
	}

	err = s.cstore.CreateContainerBundle(cont.ID(), spec, opts.RootfsPath)
	if err != nil {
		return
	}

	err = s.runtime.CreateContainer(cont.ID(), h.BundleDir())
	return
}

func (r *runtimeService) StartContainer(
	id container.ID,
) error {
	// TODO: add mutex lock to the method

	cont := r.cmap.Get(id)
	if cont == nil {
		return errors.New("container not found")
	}

	// TODO: assert cont.Status() == "created"

	if err := r.runtime.StartContainer(cont.ID()); err != nil {
		return err
	}
	return cont.SetStatus(container.StatusRunning)
}

func (r *runtimeService) StopContainer(
	id container.ID,
	timeout time.Duration,
) error {
	// TODO: add mutex lock to the method

	cont := r.cmap.Get(id)
	if cont == nil {
		return errors.New("container not found")
	}

	return r.runtime.KillContainer(cont.ID(), syscall.SIGTERM)
	// TODO: wait for `timeout` ms. If the container proc is still there
	// r.runtime.KillContainer(cont.ID(), syscall.SIGKILL)
	// wait for some default timeout. If the container proc is still there
	// kill(PID)
}

func (r *runtimeService) GetContainer(
	id container.ID,
) (*container.Container, error) {
	// TODO: add mutex lock to the method

	cont := r.cmap.Get(id)
	if cont == nil {
		return nil, errors.New("container not found")
	}

	state, err := r.runtime.ContainerState(cont.ID())
	if err != nil {
		return nil, err
	}
	status, err := container.StatusFromString(state.Status)
	if err != nil {
		return nil, err
	}
	if err := cont.SetStatus(status); err != nil {
		// TODO: warn + SetStatusForse!
	}
	// TODO: set PID
	return cont, nil
}

type ContainerOptions struct {
	Name           string
	Command        string
	Args           []string
	RootfsPath     string
	RootfsReadonly bool
}
