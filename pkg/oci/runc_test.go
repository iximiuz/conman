package oci_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/storage"
	"github.com/iximiuz/conman/pkg/testutil"
)

var cfg *config.Config

func init() {
	cfg = config.TestConfigFromFlags()
}

func Test_CreateContainer_Fail_InvalidExecutable(t *testing.T) {
	helper := NewTestHelper(t, cfg)
	defer helper.teardown()

	// Try to create an obviously broken container
	hcont, _, err := helper.createContainer("i_am_not_a_valid_executable", nil)
	if err == nil {
		t.Fatal("Expected CreateContainer() to fail!")
	}
	if !strings.Contains(err.Error(), "\\\"i_am_not_a_valid_executable\\\": executable file not found in $PATH") {
		t.Fatal("Unexpected error", err)
	}

	// Ensure shim has exited
	proc, err := findShimmyProc(hcont)
	if err != nil {
		t.Fatal(err)
	}
	if err := testutil.EnsureProcessHasTerminated(proc, 2*time.Second); err != nil {
		proc.Kill()
		t.Fatal(err)
	}
}

func Test_NonInteractive_SimpleRun(t *testing.T) {
	helper := NewTestHelper(t, cfg)
	defer helper.teardown()

	// Create container
	hcont, contPid, err := helper.createContainer("sleep", []string{"0.01"})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure shim is alive
	shimmyProc, err := findShimmyProc(hcont)
	if err != nil {
		t.Fatal(err)
	}
	if err := testutil.EnsureProcessIsAlive(shimmyProc); err != nil {
		t.Fatal(err)
	}

	// Ensure container is alive
	contProc, err := findContainerProc(hcont)
	if err != nil {
		t.Fatal(err)
	}
	if err := testutil.EnsureProcessIsAlive(contProc); err != nil {
		t.Fatal(err)
	}
	if contProc.Pid != contPid {
		t.Fatalf(
			"container pidfile (%v) != reported container pid (%v)",
			contProc.Pid, contPid,
		)
	}

	// Start container
	if err := helper.rt.StartContainer(hcont.ContainerID()); err != nil {
		t.Fatal(err)
	}

	// Ensure shim has exited
	if err := testutil.EnsureProcessHasTerminated(shimmyProc, 2*time.Second); err != nil {
		shimmyProc.Kill()
		t.Fatal(err)
	}

	// Ensure container has exited
	if err := testutil.EnsureProcessHasTerminated(contProc, 2*time.Second); err != nil {
		shimmyProc.Kill()
		t.Fatal(err)
	}

	status, err := helper.readContainerExitFile(hcont.ContainerID())
	if err != nil {
		t.Fatal(err)
	}
	if string(status) != "Exited with code 0" {
		t.Fatalf("Unexpected container termination status: %v", status)
	}

	// TODO: Check container logs
}

func Test_CreateContainer_TimeOut(t *testing.T) {
}

func Test_NonInteractive_SignalShim(t *testing.T) {
}

func Test_NonInteractive_SignalContainer(t *testing.T) {
}

type TestHelper struct {
	rt     oci.Runtime
	cstore storage.ContainerStore

	tmpDir string
}

func NewTestHelper(t *testing.T, cfg *config.Config) *TestHelper {
	tmpDir := testutil.TempDir(t)

	return &TestHelper{
		rt: oci.NewRuntime(
			cfg.ShimmyPath,
			cfg.RuntimePath,
			fsutil.EnsureExists(path.Join(tmpDir, "runc")),
			fsutil.EnsureExists(path.Join(tmpDir, "exits")),
		),
		cstore: storage.NewContainerStore(path.Join(tmpDir, "cstore")),
		tmpDir: tmpDir,
	}
}

func (h *TestHelper) createContainer(
	command string,
	args []string,
) (hcont *storage.ContainerHandle, pid int, err error) {
	contID := container.RandID()
	hcont, err = h.cstore.CreateContainer(contID, nil)
	if err != nil {
		return
	}

	spec, err := oci.NewSpec(oci.SpecOptions{
		Command:      command,
		Args:         args,
		RootPath:     hcont.RootfsDir(),
		RootReadonly: true,
	})
	if err != nil {
		return
	}

	err = h.cstore.CreateContainerBundle(
		contID,
		spec,
		testutil.DataDir("rootfs_alpine"),
	)
	if err != nil {
		return
	}

	pid, err = h.rt.CreateContainer(
		contID,
		hcont.BundleDir(),
		h.containerLogPath(contID),
		1*time.Second,
	)
	return
}

func (h *TestHelper) containerLogPath(id container.ID) string {
	return path.Join(
		fsutil.EnsureExists(path.Join(h.tmpDir, "container-logs")),
		string(id)+".log",
	)
}

func (h *TestHelper) readContainerExitFile(id container.ID) ([]byte, error) {
	return ioutil.ReadFile(path.Join(h.tmpDir, "exits", string(id)))
}

func (h *TestHelper) teardown() {
	if err := os.RemoveAll(h.tmpDir); err != nil {
		panic(err)
	}
}

func findShimmyProc(hcont *storage.ContainerHandle) (*os.Process, error) {
	return testutil.FindProcessByPidfile(
		path.Join(hcont.BundleDir(), "shimmy.pid"),
	)
}

func findContainerProc(hcont *storage.ContainerHandle) (*os.Process, error) {
	return testutil.FindProcessByPidfile(
		path.Join(hcont.BundleDir(), "container.pid"),
	)
}
