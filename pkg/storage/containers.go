package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/rollback"
)

const (
	DirAccessFailed string = "can't access container directory"
)

type ContainerStore interface {
	RootDir() string

	// CreateContainer creates container's dir in a non-volatile location
	// (it also may store some container's metadata inside).
	CreateContainer(
		container.ID,
		*rollback.Rollback,
	) (*ContainerHandle, error)

	CreateContainerBundle(
		id container.ID,
		spec oci.RuntimeSpec,
		rootfs string,
	) error

	GetContainer(container.ID) (*ContainerHandle, error)

	// Removes <container_dir>.
	DeleteContainer(container.ID) error

	// Atomically (using os.Rename) updates container's state on disk.
	// Container state is stored in <container_dir>/state.json.
	AtomicWriteContainerState(id container.ID, state []byte) error

	// Unlinks <container_dir>/state.json file effectively marking
	// the container as ready to be cleaned up.
	AtomicDeleteContainerState(container.ID) error
}

func NewContainerStore(rootdir string) ContainerStore {
	return &containerStore{
		rootdir: rootdir,
	}
}

type containerStore struct {
	rootdir string
}

func (s *containerStore) RootDir() string {
	return s.rootdir
}

func (s *containerStore) CreateContainer(
	id container.ID,
	rb *rollback.Rollback,
) (*ContainerHandle, error) {
	if rb != nil {
		rb.Add(func() { s.DeleteContainer(id) })
	}

	dir := s.containerDir(id)
	if ok, err := fsutil.Exists(dir); ok || err != nil {
		if ok {
			return nil, errors.New("container directory already exists")
		}
		return nil, errors.Wrap(err, DirAccessFailed)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, errors.Wrap(err, "can't create container directory")
	}
	return &ContainerHandle{
		containerID:  id,
		containerDir: dir,
	}, nil
}

func (s *containerStore) CreateContainerBundle(
	id container.ID,
	spec oci.RuntimeSpec,
	rootfs string,
) error {
	h, err := s.GetContainer(id)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(h.BundleDir(), 0700); err != nil {
		return errors.Wrap(err, "can't create bundle directory")
	}
	if err := fsutil.CopyDir(rootfs, h.RootfsDir()); err != nil {
		return errors.Wrap(err, "can't copy rootfs directory")
	}
	if err := ioutil.WriteFile(h.RuntimeSpecFile(), spec, 0644); err != nil {
		return errors.Wrap(err, "can't write OCI runtime spec file")
	}
	return nil
}

func (s *containerStore) GetContainer(
	id container.ID,
) (*ContainerHandle, error) {
	dir := s.containerDir(id)
	ok, err := fsutil.Exists(dir)
	if err != nil {
		return nil, errors.Wrap(err, DirAccessFailed)
	}
	if ok {
		return &ContainerHandle{
			containerID:  id,
			containerDir: dir,
		}, nil
	}
	return nil, nil
}

func (s *containerStore) DeleteContainer(id container.ID) error {
	return errors.Wrap(os.RemoveAll(s.containerDir(id)),
		"can't remove container directory")
}

func (s *containerStore) AtomicWriteContainerState(
	id container.ID,
	state []byte,
) error {
	h, err := s.GetContainer(id)
	if err != nil {
		return err
	}

	statefile := h.stateFile()
	tmpfile := statefile + ".writing"
	if err := ioutil.WriteFile(tmpfile, state, 0600); err != nil {
		return err
	}

	return os.Rename(tmpfile, statefile)
}

func (s *containerStore) AtomicDeleteContainerState(id container.ID) error {
	h, err := s.GetContainer(id)
	if err != nil {
		return err
	}
	return os.Remove(h.stateFile())
}

func (s *containerStore) containerDir(id container.ID) string {
	return path.Join(s.containersDir(), string(id))
}

func (s *containerStore) containersDir() string {
	return path.Join(s.rootdir, "containers")
}

type ContainerHandle struct {
	containerID  container.ID
	containerDir string
}

func (h *ContainerHandle) ContainerID() container.ID {
	return h.containerID
}

func (h *ContainerHandle) ContainerDir() string {
	return h.containerDir
}

func (h *ContainerHandle) BundleDir() string {
	return path.Join(h.ContainerDir(), "bundle")
}

func (h *ContainerHandle) RootfsDir() string {
	return path.Join(h.BundleDir(), "rootfs")
}

func (h *ContainerHandle) RuntimeSpecFile() string {
	return path.Join(h.BundleDir(), "config.json")
}

func (h *ContainerHandle) stateFile() string {
	return path.Join(h.ContainerDir(), "state.json")
}
