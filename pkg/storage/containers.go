package storage

import (
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/rollback"
)

const (
	DirAccessFailed string = "can't access container directory"
)

type ContainerStore interface {
	// CreateContainer creates container's dir in a non-volatile location
	// (it also may store some container's metadata inside).
	CreateContainer(*container.Container, *rollback.Rollback) error

	// CopyContainerBundle(runtime.ContainerID, path string)

	// CopyContainerSpec(runtime.ContainerID, spec)

	GetContainer(container.ID) (interface{}, error)

	DeleteContainer(container.ID) error
}

func NewContainerStore(rootdir string) ContainerStore {
	return &containerStore{
		rootdir: rootdir,
	}
}

type containerStore struct {
	sync.Mutex
	rootdir string
}

func (s *containerStore) CreateContainer(
	c *container.Container,
	rb *rollback.Rollback,
) error {
	s.Lock()
	defer s.Unlock()

	if rb != nil {
		rb.Add(func() {
			s.DeleteContainer(c.ID())
		})
	}

	dir := s.containerDir(c.ID())
	if ok, err := fsutil.Exists(dir); ok || err != nil {
		if ok {
			return errors.New("container directory already exists")
		}
		return errors.Wrap(err, DirAccessFailed)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrap(err, "can't create container directory")
	}
	return nil
}

func (s *containerStore) GetContainer(id container.ID) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	ok, err := fsutil.Exists(s.containerDir(id))
	if err != nil {
		return nil, errors.Wrap(err, DirAccessFailed)
	}
	if ok {
		return &struct{}{}, nil
	}
	return nil, nil
}

func (s *containerStore) DeleteContainer(id container.ID) error {
	s.Lock()
	defer s.Unlock()

	return errors.Wrap(os.RemoveAll(s.containerDir(id)),
		"can't remove container directory")
}

func (s *containerStore) containerDir(id container.ID) string {
	return path.Join(s.containersDir(), string(id))
}

func (s *containerStore) containersDir() string {
	return path.Join(s.rootdir, "containers")
}
