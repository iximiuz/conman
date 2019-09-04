package cri_test

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/cri"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/pkg/testutil"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestFullCycle_Simple(t *testing.T) {
	cfg, err := config.TestConfig()
	if err != nil {
		t.Fatal(err)
	}

	ociRt, teardown1 := newOciRuntime(t, cfg)
	defer teardown1()

	cstore, teardown2 := newContainerStore(t)
	defer teardown2()

	sut := cri.NewRuntimeService(ociRt, cstore)

	// (1) Create container.
	opts := cri.ContainerOptions{
		Name: "cont1",
		// Command: "/bin/sleep",
		// Args:           []string{"999"},
		Command:        "/bin/sh",
		Args:           []string{"--version"},
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
		t.Errorf("cri.StopContainer() failed.\nerr=%v\n", err)
	}

	assertContainerStatus(t, sut, contID, container.Stopped)

	// (4) RemoveContainer. TODO
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
