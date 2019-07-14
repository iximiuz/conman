package storage

import (
	"io/ioutil"
	"os"
	"path"
	// "sync"

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
		*container.Container,
		*rollback.Rollback,
	) (*ContainerHandle, error)

	CreateContainerBundle(
		id container.ID,
		spec oci.RuntimeSpec,
		rootfs string,
	) error

	GetContainer(container.ID) (*ContainerHandle, error)

	DeleteContainer(container.ID) error
}

func NewContainerStore(rootdir string) ContainerStore {
	return &containerStore{
		rootdir: rootdir,
	}
}

type containerStore struct {
	// sync.Mutex
	rootdir string
}

func (s *containerStore) RootDir() string {
	return s.rootdir
}

func (s *containerStore) CreateContainer(
	c *container.Container,
	rb *rollback.Rollback,
) (*ContainerHandle, error) {
	// s.Lock()
	// defer s.Unlock()

	if rb != nil {
		rb.Add(func() {
			s.DeleteContainer(c.ID())
		})
	}

	dir := s.containerDir(c.ID())
	if ok, err := fsutil.Exists(dir); ok || err != nil {
		if ok {
			return nil, errors.New("container directory already exists")
		}
		return nil, errors.Wrap(err, DirAccessFailed)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, errors.Wrap(err, "can't create container directory")
	}
	return &ContainerHandle{containerDir: dir}, nil
}

func (s *containerStore) CreateContainerBundle(
	id container.ID,
	spec oci.RuntimeSpec,
	rootfs string,
) error {
	// s.Lock()
	// defer s.Unlock()

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
	// s.Lock()
	// defer s.Unlock()

	dir := s.containerDir(id)
	ok, err := fsutil.Exists(dir)
	if err != nil {
		return nil, errors.Wrap(err, DirAccessFailed)
	}
	if ok {
		return &ContainerHandle{containerDir: dir}, nil
	}
	return nil, nil
}

func (s *containerStore) DeleteContainer(id container.ID) error {
	// s.Lock()
	// defer s.Unlock()

	return errors.Wrap(os.RemoveAll(s.containerDir(id)),
		"can't remove container directory")
}

func (s *containerStore) containerDir(id container.ID) string {
	return path.Join(s.containersDir(), string(id))
}

func (s *containerStore) containersDir() string {
	return path.Join(s.rootdir, "containers")
}

type ContainerHandle struct {
	containerDir string
}

func (h *ContainerHandle) ContainerDir() string {
	return h.containerDir
}

func (h *ContainerHandle) BundleDir() string {
	return path.Join(h.containerDir, "bundle")
}

func (h *ContainerHandle) RootfsDir() string {
	return path.Join(h.BundleDir(), "rootfs")
}

func (h *ContainerHandle) RuntimeSpecFile() string {
	return path.Join(h.BundleDir(), "config.json")
}
