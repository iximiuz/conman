package oci_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/fsutil"
	"github.com/iximiuz/conman/pkg/oci"
	"github.com/iximiuz/conman/pkg/shimutil"
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
	hcont, contPid, err := helper.createContainer(
		"sh",
		[]string{
			"-c",
			"echo 'stdout line'; >&2 echo 'stderr line'",
		},
	)
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

	ts, err := helper.readContainerExitFile(hcont.ContainerID())
	if err != nil {
		t.Fatal(err)
	}
	if ts.IsSignaled() == true {
		t.Fatalf("Unexpected container termination status: %+v", ts)
	}
	if ts.ExitCode() != 0 {
		t.Fatalf("Unexpected container termination status: %+v", ts)
	}

	// Validate container logs
	containerLog, err := helper.readContainerLog(hcont.ContainerID())
	if err != nil {
		t.Fatal(err)
	}
	if containerLog.stdout[0].message != "stdout line" {
		t.Fatalf("Unexpected container log: %+v", containerLog)
	}
	if containerLog.stderr[0].message != "stderr line" {
		t.Fatalf("Unexpected container log: %+v", containerLog)
	}
}

func Test_CreateContainer_TimeOut(t *testing.T) {
}

func Test_NonInteractive_SignalShim(t *testing.T) {
}

func Test_NonInteractive_SignalContainer(t *testing.T) {
}

func Test_CreateContainer_ContainerShellExitsWithError(t *testing.T) {
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
		h.containerExitPath(contID),
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

func (h *TestHelper) containerExitPath(id container.ID) string {
	return path.Join(
		fsutil.EnsureExists(path.Join(h.tmpDir, "exits")),
		string(id),
	)
}

type ContainerLogRecord struct {
	time    string
	message string
}

type ContainerLog struct {
	stdout []ContainerLogRecord
	stderr []ContainerLogRecord
}

func (h *TestHelper) readContainerLog(id container.ID) (ContainerLog, error) {
	log := ContainerLog{}
	bytes, err := ioutil.ReadFile(h.containerLogPath(id))
	if err != nil {
		return log, err
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) != 3 || (parts[1] != "stdout" && parts[1] != "stderr") {
			return ContainerLog{}, errors.Errorf("Malformed log record [%s]", line)
		}

		if parts[1] == "stdout" {
			log.stdout = append(log.stdout, ContainerLogRecord{time: parts[0], message: parts[2]})
		}
		if parts[1] == "stderr" {
			log.stderr = append(log.stderr, ContainerLogRecord{time: parts[0], message: parts[2]})
		}
	}

	return log, nil
}

func (h *TestHelper) readContainerExitFile(id container.ID) (*shimutil.TerminationStatus, error) {
	bytes, err := ioutil.ReadFile(h.containerExitPath(id))
	if err != nil {
		return nil, err
	}
	return shimutil.ParseExitFile(bytes)
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
