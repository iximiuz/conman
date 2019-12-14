package shimmy_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/iximiuz/conman/config"
	"github.com/iximiuz/conman/pkg/testutil"
	"github.com/iximiuz/conman/pkg/timeutil"
)

var cfg *config.Config

func init() {
	cfg = config.TestConfigFromFlags()
}

func TestAbnormalRuntimeTermination(t *testing.T) {
	tmpdir := testutil.TempDir(t)
	defer os.RemoveAll(tmpdir)

	pidfile := path.Join(tmpdir, "shimmy.pid")

	cmd := exec.Command(
		cfg.ShimmyPath,
		"--shimmy-pidfile", pidfile,
		"--runtime", cfg.RuntimePath,
		"--runtime-arg", "foobar=123",
		"--bundle", "/not/used/folder",
		"--container-id", "<not-used-id>",
		"--container-pidfile", "/not/used/file.pid",
		"--container-logfile", "/not/used/logfile",
		"--container-exitfile", "/not/used/exitfile",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	rd, wr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer rd.Close()
	defer wr.Close()

	cmd.ExtraFiles = append(cmd.ExtraFiles, wr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("_OCI_SYNCPIPE=%d", 3))

	type Report struct {
		Kind   string `json:"kind"`
		Status string `json:"status"`
		Stderr string `json:"stderr"`
	}

	if err = timeutil.WithTimeout(3*time.Second, func() error {
		if err := cmd.Run(); err != nil {
			return err
		}

		wr.Close()

		bytes, err := ioutil.ReadAll(rd)
		if err != nil {
			return err
		}

		rd.Close()

		report := Report{}
		if err := json.Unmarshal(bytes, &report); err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("Failed to decode report string [%v]. Raw [%v].",
					string(bytes), bytes),
			)
		}

		if report.Kind != "runtime_abnormal_termination" {
			return errors.Errorf("Unexpected report kind: %v", string(bytes))
		}
		if report.Status != "Runtime Exited with code 3." {
			return errors.Errorf("Unexpected report status: %v", string(bytes))
		}
		if !strings.Contains(report.Stderr, "No help topic for 'foobar=123'") {
			return errors.Errorf("Unexpected report stderr: %v", string(bytes))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	proc, err := testutil.FindProcessByPidfile(pidfile)
	if err != nil {
		t.Fatal(err)
	}

	if err := testutil.EnsureProcessHasTerminated(proc, 2*time.Second); err != nil {
		proc.Kill()
		t.Fatal(err)
	}
}
