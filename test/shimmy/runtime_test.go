package shimmy_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultShimmyExe = "/usr/local/bin/shimmy"
)

var (
	// path to shimmy executable, eg. /usr/local/bin/shimmy
	shimmyExe string
)

func init() {
	flag.StringVar(&shimmyExe, "shimmy", defaultShimmyExe, "Path to shimmy executable file")
	flag.Parse()
}

func withTimeout(d time.Duration, fn func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	ch := make(chan error, 1)

	func() {
		ch <- fn()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "Timed out")
	}
}

func TestAbnormalRuntimeTermination(t *testing.T) {
	cmd := exec.Command(shimmyExe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	rw, wr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer rw.Close()
	defer wr.Close()

	cmd.ExtraFiles = append(cmd.ExtraFiles, wr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("_OCI_SYNCPIPE=%d", 3))

	type Report struct {
		Kind   string `json:"kind"`
		Status string `json:"status"`
		Stderr string `json:"stderr"`
	}

	if err = withTimeout(3*time.Second, func() error {
		if err := cmd.Run(); err != nil {
			return err
		}

		wr.Close()

		bytes, err := ioutil.ReadAll(rw)
		if err != nil {
			return err
		}

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
		if report.Status != "Runtime Exited with code 1." {
			return errors.Errorf("Unexpected report status: %v", string(bytes))
		}
		if !strings.Contains(report.Stderr, "flag provided but not defined: -foobar") {
			return errors.Errorf("Unexpected report stderr: %v", string(bytes))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
