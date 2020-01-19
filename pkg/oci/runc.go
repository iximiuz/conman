package oci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/timeutil"
)

// implementes oci.Runtime interface
type runcRuntime struct {
	// path to shimmy executable, eg. /usr/local/bin/shimmy
	shimmyPath string

	// path to runc executable, eg. /usr/bin/runc
	runtimePath string

	// dir to store container state (on tmpfs), eg. /run/runc/
	rootPath string
}

func NewRuntime(
	shimmyPath string,
	runtimePath string,
	rootPath string,
) Runtime {
	return &runcRuntime{
		shimmyPath:  shimmyPath,
		runtimePath: runtimePath,
		rootPath:    rootPath,
	}
}

func (r *runcRuntime) CreateContainer(
	id container.ID,
	bundleDir string,
	logfile string,
	exitfile string,
	attachfile string,
	stdin bool,
	stdinOnce bool,
	timeout time.Duration,
) (pid int, err error) {
	cmd := exec.Command(
		r.shimmyPath,
		"--shimmy-pidfile", path.Join(bundleDir, "shimmy.pid"),
		"--shimmy-log-level", strings.ToUpper(logrus.GetLevel().String()),
		"--runtime", r.runtimePath,
		"--runtime-arg", fmt.Sprintf("'--root=%s'", r.rootPath),
		"--bundle", bundleDir,
		"--container-id", string(id),
		"--container-pidfile", path.Join(bundleDir, "container.pid"),
		"--container-logfile", logfile,
		"--container-exitfile", exitfile,
		"--container-attachfile", attachfile,
	)
	if stdin {
		cmd.Args = append(cmd.Args, "--stdin")
	}
	if stdinOnce {
		cmd.Args = append(cmd.Args, "--stdin-once")
	}

	syncpipeRead, syncpipeWrite, err := os.Pipe()
	if err != nil {
		return 0, err
	}
	defer syncpipeRead.Close()
	defer syncpipeWrite.Close()

	cmd.ExtraFiles = append(cmd.ExtraFiles, syncpipeWrite)
	cmd.Args = append(
		cmd.Args,
		// 0,1, and 2 are STDIO streams
		"--syncpipe-fd", strconv.Itoa(2+len(cmd.ExtraFiles)),
	)

	// We expect shimmy execution to be almost instant, because its
	// main process just validates the input parameters, forks the
	// shim process, saves its PID on disk, and then exits.
	if _, err := runCommand(cmd); err != nil {
		return 0, err
	}

	syncpipeWrite.Close()

	type Report struct {
		Kind   string `json:"kind"`
		Status string `json:"status"`
		Stderr string `json:"stderr"`
		Pid    int    `json:"pid"`
	}

	return pid, timeutil.WithTimeout(timeout, func() error {
		bytes, err := ioutil.ReadAll(syncpipeRead)
		if err != nil {
			return err
		}
		syncpipeRead.Close()

		report := Report{}
		if err := json.Unmarshal(bytes, &report); err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("Failed to decode report string [%v]. Raw [%v].",
					string(bytes), bytes),
			)
		}

		if report.Kind == "container_pid" && report.Pid > 0 {
			pid = report.Pid
			return nil
		}
		return errors.Errorf("%+v", report)
	})
}

func (r *runcRuntime) StartContainer(id container.ID) error {
	cmd := exec.Command(
		r.runtimePath,
		"--root", r.rootPath,
		"start", string(id),
	)
	_, err := runCommand(cmd)
	return err
}

func (r *runcRuntime) KillContainer(id container.ID, sig os.Signal) error {
	sigstr, err := sigStr(sig)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		r.runtimePath,
		"--root", r.rootPath,
		"kill",
		string(id),
		sigstr,
	)
	_, err = runCommand(cmd)
	return err
}

func (r *runcRuntime) DeleteContainer(id container.ID) error {
	cmd := exec.Command(
		r.runtimePath,
		"--root", r.rootPath,
		"delete",
		string(id),
	)
	_, err := runCommand(cmd)
	return err
}

func (r *runcRuntime) ContainerState(id container.ID) (StateResp, error) {
	cmd := exec.Command(
		r.runtimePath,
		"--root", r.rootPath,
		"state",
		string(id),
	)
	output, err := runCommand(cmd)
	if err != nil {
		return StateResp{}, err
	}

	resp := StateResp{}
	return resp, json.Unmarshal(output, &resp)
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.Output()
	debugLog(cmd, output, err)
	return output, wrappedError(err)
}

func debugLog(c *exec.Cmd, stdout []byte, err error) {
	stderr := []byte{}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = ee.Stderr
		}
	}
	logrus.WithFields(logrus.Fields{
		"stdout": string(stdout),
		"stderr": string(stderr),
		"error":  err,
	}).Debugf("exec %s", strings.Join(c.Args, " "))
}

func wrappedError(err error) error {
	if err == nil {
		return nil
	}

	msg := "OCI runtime (runc) execution failed"
	if ee, ok := err.(*exec.ExitError); ok {
		msg = fmt.Sprintf("%v, stderr=[%v]", msg, string(ee.Stderr))
	}
	return errors.Wrap(err, msg)
}
