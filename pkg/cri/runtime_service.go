package cri

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/kubernetes/pkg/kubelet/cri/streaming"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/rollback"
	"github.com/iximiuz/conman/pkg/shimutil"
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

	// Removes container from both conman and runc storages.
	// If container has not been stopped yet, a force flag
	// must be set. If container has already been removed, no
	// error returned (i.e. idempotent behavior).
	RemoveContainer(container.ID) error

	ListContainers() ([]*container.Container, error)

	// GetContainer returns the container doing a state request
	// from the OCI runtime if applicable.
	GetContainer(container.ID) (*container.Container, error)

	// TODO: UpdateContainerResources
	// TODO: ReopenContainerLog

	streaming.Runtime

	// TODO: RunPodSandbox
	// TODO: StopPodSandbox
	// TODO: RemovePodSandbox
	// TODO: PodSandboxStatus
	// TODO: ListPodSandbox
}

// runtimeService implements RuntimeService interface.
// Some design considerations:
//   - runtimeService methods are thread-safe. There is a common locker
//     on runtimeService instance protecting from concurrent container
//     modifications. Given this lock, dependencies like container.Map,
//     storage.ContainerStore can be simplified and omit their own locking.
//   - runtimeService tracks container states on its own. It uses ContainerStore
//     to write a JSON-serialized container state inside of container base dir.
//     Since atomic write of the state and runc execution is not possible,
//     the state modification happens first (optimistic approach). Then a runc
//     command follows. In case of the runc error, the state should be rolled
//     back. However, during a cascading failure, the state saved in the container
//     dir and the state of the container in accordance with runc might diverge.
//     The state restoring logic should try to fix the introduced discrepancy.
//   - ContainerStore is the only source of truth. Only containers tracked by
//     the store are managed by conman. I.e. if someone uses the same runc config
//     to create extra containers, the change will not be visible to conman.
//     At the same time modification of the managed containers by running runc
//     manually (or by any other means) will introduce inconsistency in the conman
//     tracked state and the actual state of the containers.
type runtimeService struct {
	sync.Mutex

	runtime   oci.Runtime
	cstore    storage.ContainerStore
	logDir    string
	exitDir   string
	attachDir string

	cmap *container.Map
}

func NewRuntimeService(
	runtime oci.Runtime,
	cstore storage.ContainerStore,
	logDir string,
	exitDir string,
	attachDir string,
) (RuntimeService, error) {
	rs := &runtimeService{
		runtime:   runtime,
		cstore:    cstore,
		logDir:    logDir,
		exitDir:   exitDir,
		attachDir: attachDir,
		cmap:      container.NewMap(),
	}
	if err := rs.restore(); err != nil {
		return nil, err
	}
	return rs, nil
}

func (rs *runtimeService) CreateContainer(
	opts ContainerOptions,
) (cont *container.Container, err error) {
	rs.Lock()
	defer rs.Unlock()

	rb := rollback.New()
	defer func() { _ = err == nil || rb.Execute() }()

	contID := container.RandID()
	cont, err = container.New(
		contID,
		opts.Name,
		rs.containerLogFile(contID),
	)
	if err != nil {
		return
	}

	if err = rs.cmap.Add(cont, rb); err != nil {
		return
	}

	hcont, err := rs.cstore.CreateContainer(cont.ID(), rb)
	if err != nil {
		return
	}

	spec, err := oci.NewSpec(oci.SpecOptions{
		Command:      opts.Command,
		Args:         opts.Args,
		RootPath:     hcont.RootfsDir(),
		RootReadonly: opts.RootfsReadonly,
	})
	if err != nil {
		return
	}

	err = rs.cstore.CreateContainerBundle(cont.ID(), spec, opts.RootfsPath)
	if err != nil {
		return
	}

	err = rs.optimisticChangeContainerStatus(cont, container.Created)
	if err != nil {
		return
	}

	_, err = rs.runtime.CreateContainer(
		cont.ID(),
		hcont.BundleDir(),
		cont.LogPath(),
		rs.containerExitFile(cont.ID()),
		rs.containerAttachFile(cont.ID()),
		opts.Stdin,
		opts.StdinOnce,
		10*time.Second,
	)
	if err != nil {
		return
	}

	err = cont.SetCreatedAt(time.Now())
	return
}

func (rs *runtimeService) StartContainer(
	id container.ID,
) error {
	rs.Lock()
	defer rs.Unlock()

	cont := rs.cmap.Get(id)
	if cont == nil {
		return errors.New("container not found")
	}
	if err := assertStatus(cont.Status(), container.Created); err != nil {
		return err
	}

	if err := rs.optimisticChangeContainerStatus(cont, container.Running); err != nil {
		return err
	}

	if err := rs.runtime.StartContainer(cont.ID()); err != nil {
		return err
	}

	if err := rs.waitContainerStartedNoLock(id); err != nil {
		return nil
	}

	return cont.SetStartedAt(time.Now())
}

func (rs *runtimeService) StopContainer(
	id container.ID,
	timeout time.Duration,
) error {
	rs.Lock()
	defer rs.Unlock()

	cont := rs.cmap.Get(id)
	if cont == nil {
		return errors.New("container not found")
	}
	if err := assertStatus(
		cont.Status(), container.Created, container.Running); err != nil {
		return err
	}

	// TODO: impl a PROPPER ALGO. Wait for `timeout` ms. If the container proc
	// is still there rs.runtime.KillContainer(cont.ID(), syscall.SIGKILL)
	// wait for some default timeout. If the container proc is still there
	// os.kill(PID).

	// TODO: test for this logic!

	if err := rs.optimisticChangeContainerStatus(cont, container.Stopped); err != nil {
		return err
	}

	if err := rs.runtime.KillContainer(cont.ID(), syscall.SIGTERM); err != nil {
		return err
	}

	delays := []time.Duration{
		250 * time.Millisecond,
		250 * time.Millisecond,
	}
	for _, d := range delays {
		time.Sleep(d)
		cont, err := rs.getContainerNoLock(id)
		if err != nil {
			return err
		}
		if cont.Status() == container.Stopped {
			return nil
		}
	}

	if err := rs.runtime.KillContainer(cont.ID(), syscall.SIGKILL); err != nil {
		return err
	}
	for _, d := range delays {
		time.Sleep(d)
		cont, err := rs.getContainerNoLock(id)
		if err != nil {
			return err
		}
		if cont.Status() == container.Stopped {
			return nil
		}
	}

	return errors.New("Cannot kill container. TODO: use os.kill() to force kill")
}

func (rs *runtimeService) RemoveContainer(id container.ID) error {
	rs.Lock()
	defer rs.Unlock()

	cont := rs.cmap.Get(id)
	if cont == nil {
		return nil
	}

	// Atomically mark container removed
	if err := rs.cstore.ContainerStateDeleteAtomic(id); err != nil {
		return err
	}

	// Initiate actual removal
	if err := rs.runtime.DeleteContainer(cont.ID()); err != nil {
		return err
	}

	// Cleanup leftovers
	rs.cmap.Del(id)
	return rs.cstore.DeleteContainer(id)
}

func (rs *runtimeService) ListContainers() ([]*container.Container, error) {
	rs.Lock()
	defer rs.Unlock()

	var cs []*container.Container
	for _, c := range rs.cmap.All() {
		c, err := rs.getContainerNoLock(c.ID())
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}

	sort.SliceStable(cs, func(i, j int) bool {
		iat := cs[i].CreatedAtNano()
		jat := cs[j].CreatedAtNano()
		if iat == jat {
			return cs[i].ID() < cs[j].ID()
		}
		return iat < jat
	})

	return cs, nil
}

func (rs *runtimeService) GetContainer(
	id container.ID,
) (*container.Container, error) {
	rs.Lock()
	defer rs.Unlock()
	return rs.getContainerNoLock(id)
}

func (rs *runtimeService) getContainerNoLock(
	id container.ID,
) (*container.Container, error) {
	cont := rs.cmap.Get(id)
	if cont == nil {
		return nil, errors.New("container not found")
	}

	// Request container state
	state, err := rs.runtime.ContainerState(cont.ID())
	if err != nil {
		return nil, err
	}

	// Set container status from state
	status, err := container.StatusFromString(state.Status)
	if err != nil {
		return nil, err
	}
	cont.SetStatus(status)

	// Set container exit code if applicable
	if cont.Status() == container.Stopped {
		ts, err := rs.parseContainerExitFile(id)
		if err != nil {
			return nil, err
		}

		cont.SetFinishedAt(ts.At())

		if ts.IsSignaled() {
			cont.SetExitCode(127 + ts.Signal())
		} else {
			cont.SetExitCode(ts.ExitCode())
		}
	}

	blob, err := cont.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := rs.cstore.ContainerStateWriteAtomic(id, blob); err != nil {
		return nil, err
	}

	return cont, nil
}

func (rs *runtimeService) restore() error {
	rs.Lock()
	defer rs.Unlock()

	hconts, err := rs.cstore.FindContainers()
	if err != nil {
		return err
	}

	purgeBrokenContainer := func(id container.ID) {
		rs.cmap.Del(id)
		if err := rs.cstore.DeleteContainer(id); err != nil {
			logrus.WithError(err).Warn("failed to purge broken container")
		}
	}

	for _, h := range hconts {
		blob, err := rs.cstore.ContainerStateRead(h.ContainerID())
		if err != nil {
			logrus.WithError(err).Warn("failed to read container state")
			purgeBrokenContainer(h.ContainerID())
			continue
		}

		cont := &container.Container{}
		if err := cont.UnmarshalJSON(blob); err != nil {
			logrus.WithError(err).Warn("failed to unmarshal container state")
			continue
		}

		if err := rs.cmap.Add(cont, nil); err != nil {
			logrus.WithError(err).Warn("failed to in-memory store container")
			continue
		}

		cont, err = rs.getContainerNoLock(h.ContainerID())
		if err != nil {
			logrus.WithError(err).Warn("failed to update container state")
			purgeBrokenContainer(h.ContainerID())
			continue
		}
	}

	return nil
}

func (rs *runtimeService) waitContainerStartedNoLock(id container.ID) error {
	delays := []time.Duration{
		250 * time.Millisecond,
		250 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
	}
	status := container.Unknown
	for _, d := range delays {
		time.Sleep(d)
		cont, err := rs.getContainerNoLock(id)
		status = cont.Status()
		if err != nil {
			return err
		}
		if status == container.Running {
			return nil
		}
		if status != container.Created {
			break
		}
	}

	// TODO: handle case with containers exiting quickly
	return errors.New(
		fmt.Sprintf("Failed to start container; status=%v.", status))
}

func (rs *runtimeService) optimisticChangeContainerStatus(
	c *container.Container,
	s container.Status,
) error {
	c.SetStatus(s)
	blob, err := c.MarshalJSON()
	if err != nil {
		return err
	}
	return rs.cstore.ContainerStateWriteAtomic(c.ID(), blob)
}

func (rs *runtimeService) containerLogFile(id container.ID) string {
	return path.Join(rs.logDir, string(id)+".log")
}

func (rs *runtimeService) containerAttachFile(id container.ID) string {
	return path.Join(rs.attachDir, string(id))
}

func (rs *runtimeService) containerExitFile(id container.ID) string {
	return path.Join(rs.exitDir, string(id))
}

func (rs *runtimeService) parseContainerExitFile(id container.ID) (*shimutil.TerminationStatus, error) {
	bytes, err := ioutil.ReadFile(rs.containerExitFile(id))
	if err != nil {
		return nil, errors.Wrap(err, "container exit file parsing failed")
	}
	return shimutil.ParseExitFile(bytes)
}

type ContainerOptions struct {
	Name           string
	Command        string
	Args           []string
	RootfsPath     string
	RootfsReadonly bool
	Stdin          bool
	StdinOnce      bool
}

func assertStatus(actual container.Status, expected ...container.Status) error {
	for _, e := range expected {
		if actual == e {
			return nil
		}
	}
	return errors.New(
		fmt.Sprintf("Wrong container status \"%v\". Expected one of=%v",
			actual, expected))
}
