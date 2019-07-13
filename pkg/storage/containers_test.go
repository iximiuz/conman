package storage

import (
	"os"
	"testing"

	"github.com/iximiuz/conman/pkg/testutil"
)

func TestCreateContainer(t *testing.T) {
	dir := testutil.TempDir(t)
	defer os.RemoveAll(dir)
	s := NewContainerStore(dir)

	c := testutil.NewContainer()
	if err := s.CreateContainer(c, nil); err != nil {
		t.Fatal("ContainerStore cannot create container")
	}
}

func TestGetContainer(t *testing.T) {
	dir := testutil.TempDir(t)
	defer os.RemoveAll(dir)
	s := NewContainerStore(dir)

	c := testutil.NewContainer()
	if err := s.CreateContainer(c, nil); err != nil {
		t.Fatal("ContainerStore cannot create container")
	}

	h, err := s.GetContainer(c.ID())
	if err != nil {
		t.Fatal("GetContainer() failed with error", err)
	}
	if h == nil {
		t.Fatal("ContainerStore cannot find container")
	}
}
