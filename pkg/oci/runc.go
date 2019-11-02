package oci

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/iximiuz/conman/pkg/container"
)

// implementes oci.Runtime interface
type runcRuntime struct {
	// path to runc executable, eg. /usr/local/bin/runc
	exePath string

	// dir to store container state (on tmpfs), eg. /run/runc/
	rootPath string
}

func NewRuntime(exePath, rootPath string) Runtime {
	return &runcRuntime{
		exePath:  exePath,
		rootPath: rootPath,
	}
}

func (r *runcRuntime) CreateContainer(
	id container.ID,
	containerDir string,
	bundleDir string,
) error {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"--log", path.Join(containerDir, "runc.log"),
		"create",
		"--bundle", bundleDir,
		string(id),
	)

	// Cannot use cmd.Output() here 'cause runc forks a child process
	// and its standard streams are still connected to the runc
	// streams. So, even though the parent (i.e. runc) process terminates
	// conman will wait till the output (stdout and stderr) streams have
	// been closed.

	// Do a dirty trick here - in order to have at least some visibility
	// on runc stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	debugLog(cmd, nil, err)
	return wrappedError(err)
}

func (r *runcRuntime) StartContainer(
	id container.ID,
	containerDir string,
	bundleDir string,
) error {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"--log", path.Join(containerDir, "runc.log"),
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
		r.exePath,
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
		r.exePath,
		"--root", r.rootPath,
		"delete",
		string(id),
	)
	_, err := runCommand(cmd)
	return err
}

func (r *runcRuntime) ContainerState(id container.ID) (StateResp, error) {
	cmd := exec.Command(
		r.exePath,
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
