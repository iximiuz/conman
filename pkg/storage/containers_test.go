package storage

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/testutil"
)

func TestCreateContainer(t *testing.T) {
	dir := testutil.TempDir(t)
	defer os.RemoveAll(dir)
	s := NewContainerStore(dir)

	c := testutil.NewContainer()
	h, err := s.CreateContainer(c, nil)
	if err != nil {
		t.Fatal("ContainerStore cannot create container")
	}
	if h == nil {
		t.Fatal("ContainerStore returned no handle")
	}
}

func TestGetContainer(t *testing.T) {
	s, c := storeWithContainer(t)
	defer os.RemoveAll(s.RootDir())

	h, err := s.GetContainer(c.ID())
	if err != nil {
		t.Fatal("GetContainer() failed with error", err)
	}
	if h == nil {
		t.Fatal("ContainerStore cannot find container")
	}
}

func TestCreateContainerBundle(t *testing.T) {
	s, c := storeWithContainer(t)
	defer os.RemoveAll(s.RootDir())

	rootfs := makeRootfs(t)
	defer os.RemoveAll(rootfs)

	spec := oci.RuntimeSpec("{}")
	err := s.CreateContainerBundle(c.ID(), spec, rootfs)
	if err != nil {
		t.Fatal("ContainerStore failed to create bundle", err)
	}
}

func storeWithContainer(t *testing.T) (ContainerStore, *container.Container) {
	dir := testutil.TempDir(t)
	s := NewContainerStore(dir)

	c := testutil.NewContainer()
	h, err := s.CreateContainer(c, nil)
	if err != nil {
		t.Fatal("ContainerStore cannot create container")
	}
	if h == nil {
		t.Fatal("ContainerStore returned no handle")
	}
	return s, c
}

func makeRootfs(t *testing.T) string {
	rootfs := testutil.TempDir(t)
	must(os.MkdirAll(path.Join(rootfs, "qux"), 0700))
	must(ioutil.WriteFile(path.Join(rootfs, "a.txt"), []byte("foo"), 0644))
	must(ioutil.WriteFile(path.Join(rootfs, "qux", "b.txt"), []byte("bar"), 0644))
	return rootfs
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
