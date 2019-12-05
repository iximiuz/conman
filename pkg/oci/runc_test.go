package oci_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/pkg/testutil"
)

var cfg *config.Config

func init() {
	cfg = config.TestConfigFromFlags()
}

func Test_CreateContainer_Fail(t *testing.T) {
	rt, teardown1 := newOciRuntime(t, cfg)
	defer teardown1()

	cstore, teardown2 := newContainerStore(t)
	defer teardown2()

	contID := container.RandID()
	hcont, err := cstore.CreateContainer(contID, nil)
	if err != nil {
		t.Fatal(err)
	}

	spec, err := oci.NewSpec(oci.SpecOptions{
		Command:      "i_am_not_a_valid_executable",
		RootPath:     hcont.RootfsDir(),
		RootReadonly: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := cstore.CreateContainerBundle(
		contID,
		spec,
		testutil.DataDir("rootfs_alpine"),
	); err != nil {
		t.Fatal(err)
	}

	err = rt.CreateContainer(
		contID,
		hcont.ContainerDir(),
		hcont.BundleDir(),
		1*time.Second,
	)
	if err == nil {
		t.Fatal("Expected CreateContainer() to fail!")
	}
	if !strings.Contains(err.Error(), "\\\"i_am_not_a_valid_executable\\\": executable file not found in $PATH") {
		t.Fatal("Unexpected error", err)
	}

	// TODO: assert pid file exists
	// TODO: assert no process with such pid exists
}

func Test_CreateContainer_TimeOut(t *testing.T) {
}

func Test_NonInteractive_SimpleRun(t *testing.T) {
}

func Test_NonInteractive_SignalShim(t *testing.T) {
}

func Test_NonInteractive_SignalContainer(t *testing.T) {
}

func Test_NonInteractive_KillShim(t *testing.T) {
}

func newOciRuntime(
	t *testing.T,
	cfg *config.Config,
) (oci.Runtime, func()) {
	root := testutil.TempDir(t)
	return oci.NewRuntime(
		cfg.ShimmyPath,
		cfg.RuntimePath,
		root,
	), func() { os.RemoveAll(root) }
}

func newContainerStore(
	t *testing.T,
) (storage.ContainerStore, func()) {
	root := testutil.TempDir(t)
	return storage.NewContainerStore(root), func() { os.RemoveAll(root) }
}
