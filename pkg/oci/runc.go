package oci

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strings"

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
	err := cmd.Run()
	debugLog(cmd, nil, err)
	return err
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
	// Cannot use cmd.Output() here 'cause runc forks a child process
	// and its standard streams are still connected to the runc
	// streams. So, even though the parent (i.e. runc) process terminates
	// conman will wait till the output (stdout and stderr) streams have
	// been closed.
	err := cmd.Run()
	debugLog(cmd, nil, err)
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
	out, err := cmd.Output()
	debugLog(cmd, out, err)
	return err
}

func (r *runcRuntime) DeleteContainer(id container.ID) error {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"delete",
		string(id),
	)
	out, err := cmd.Output()
	debugLog(cmd, out, err)
	return err
}

func (r *runcRuntime) ContainerState(id container.ID) (StateResp, error) {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"state",
		string(id),
	)
	out, err := cmd.Output()
	debugLog(cmd, out, err)
	if err != nil {
		return StateResp{}, err
	}

	resp := StateResp{}
	return resp, json.Unmarshal(out, &resp)
}

func debugLog(c *exec.Cmd, stdout []byte, err error) {
	logrus.WithFields(logrus.Fields{
		"stdout": string(stdout),
		"error":  err,
	}).Debugf("exec %s", strings.Join(c.Args, " "))
}
