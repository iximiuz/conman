package cri_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/cri"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/pkg/testutil"
)

var cfg *config.Config

func init() {
	cfg = config.TestConfigFromFlags()
}

func Test_NonInteractive_FullCycle_Simple(t *testing.T) {
	ociRt, teardown1 := newOciRuntime(t, cfg)
	defer teardown1()

	cstore, teardown2 := newContainerStore(t)
	defer teardown2()

	logdir := testutil.TempDir(t)
	defer os.RemoveAll(logdir)

	sut, err := cri.NewRuntimeService(ociRt, cstore, logdir)
	if err != nil {
		t.Fatal(err)
	}

	// (1) Create container.
	opts := cri.ContainerOptions{
		Name:           "cont1",
		Command:        "/bin/sleep",
		Args:           []string{"999"},
		RootfsPath:     testutil.DataDir("rootfs_alpine"),
		RootfsReadonly: true,
	}
	cont, err := sut.CreateContainer(opts)
	if err != nil {
		t.Fatalf("cri.CreateContainer() failed.\nerr=%v\nargs=%+v\n", err, opts)
	}
	contID := cont.ID()
	defer sut.StopContainer(contID, 500*time.Millisecond)

	assertContainerStatus(t, sut, contID, container.Created)

	// (2) Start container.
	err = sut.StartContainer(contID)
	if err != nil {
		t.Fatalf("cri.StartContainer() failed.\nerr=%v\n", err)
	}

	assertContainerStatus(t, sut, contID, container.Running)

	// (3) Stop container.
	err = sut.StopContainer(contID, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("cri.StopContainer() failed.\nerr=%v\n", err)
	}

	assertContainerStatus(t, sut, contID, container.Stopped)

	// (4) RemoveContainer.
	err = sut.RemoveContainer(contID)
	if err != nil {
		t.Fatalf("cri.RemoveContainer() failed.\nerr=%v\n", err)
	}

	_, err = sut.GetContainer(contID)
	if err == nil || err.Error() != "container not found" {
		t.Fatalf("RemoveContainer() did not remove container.\nerr=%v\n", err)
	}
}

func newOciRuntime(
	t *testing.T,
	cfg *config.Config,
) (oci.Runtime, func()) {
	tmp := testutil.TempDir(t)
	return oci.NewRuntime(
		cfg.ShimmyPath,
		cfg.RuntimePath,
		fsutil.EnsureExists(path.Join(tmp, "runc")),
		fsutil.EnsureExists(path.Join(tmp, "exits")),
	), func() { os.RemoveAll(tmp) }
}

func newContainerStore(
	t *testing.T,
) (storage.ContainerStore, func()) {
	root := testutil.TempDir(t)
	return storage.NewContainerStore(root), func() { os.RemoveAll(root) }
}

func assertContainerStatus(
	t *testing.T,
	sut cri.RuntimeService,
	id container.ID,
	expected container.Status,
) {
	cont, err := sut.GetContainer(id)
	if err != nil {
		t.Fatalf("cri.ContainerStatus() failed.\nerr=%v\n", err)
	}
	actual := cont.Status()
	if expected != actual {
		t.Fatalf("status is %v, expected status %v\n", actual, expected)
	}
}
