package storage

import (
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/rollback"
)

type ContainerStore interface {
	// CreateContainer creates container dir in non-volatile dir
	// (it also may store some metadata aside).
	CreateContainer(*container.Container, *rollback.Rollback) error

	// CopyContainerBundle(runtime.ContainerID, path string)

	// CopyContainerSpec(runtime.ContainerID, spec)

	// GetContainer(runtime.ContainerID) C1

	DeleteContainer(container.ID) error
}

func New(rootdir string) ContainerStore {
	return &containerStore{
		rootdir: rootdir,
	}
}

type containerStore struct {
	rootdir string
}

func (s *containerStore) CreateContainer(
	c *container.Container,
	rb *rollback.Rollback,
) error {
	if rb != nil {
		rb.Add(func() {
			s.DeleteContainer(c.ID())
		})
	}
	return nil
}

func (s *containerStore) DeleteContainer(_id container.ID) error {
	return nil
}
