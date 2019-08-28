package cri_test

import (
	"os"
	"testing"
	"time"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/cri"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/pkg/testutil"
)

func TestCreateContainer(t *testing.T) {
	cfg, err := config.TestConfig()
	if err != nil {
		t.Fatal(err)
	}

	ociRt, teardown1 := newOciRuntime(t, cfg)
	defer teardown1()

	cstore, teardown2 := newContainerStore(t)
	defer teardown2()

	sut := cri.NewRuntimeService(ociRt, cstore)

	opts := cri.ContainerOptions{
		Name:           "cont1",
		Command:        "/bin/sh",
		RootfsPath:     testutil.DataDir("rootfs_alpine"),
		RootfsReadonly: true,
	}
	cont, err := sut.CreateContainer(opts)
	if err != nil {
		t.Fatalf("cri.CreateContainer() failed.\nerr=%v\nargs=%+v\n", err, opts)
	}

	err = sut.StopContainer(cont.ID(), 500*time.Millisecond)
	if err != nil {
		t.Fatalf("cri.StopContainer() failed.\nerr=%v\n", err)
	}
}

func Test1(t *testing.T) {
}

func newOciRuntime(
	t *testing.T,
	cfg *config.Config,
) (oci.Runtime, func()) {
	root := testutil.TempDir(t)
	return oci.NewRuntime(cfg.RuntimePath, root), func() { os.RemoveAll(root) }
}

func newContainerStore(
	t *testing.T,
) (storage.ContainerStore, func()) {
	root := testutil.TempDir(t)
	return storage.NewContainerStore(root), func() { os.RemoveAll(root) }
}
