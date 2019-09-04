package oci

import (
	"encoding/json"
	"os"
	"os/exec"
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

func (r *runcRuntime) CreateContainer(id container.ID, bundle string) error {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"--debug",
		"--log", "/var/log/runc1.log",
		"create",
		"--bundle", bundle,
		string(id),
	)
	err := cmd.Run()
	debugLogCmd(cmd, nil, err)
	return err
}

func (r *runcRuntime) StartContainer(
	id container.ID,
	bundleDir string,
) error {
	cmd := exec.Command(
		r.exePath,
		"--debug",
		"--log", "/var/log/runc1.log",
		"--root", r.rootPath,
		"start", string(id),
	)
	cmd.Dir = bundleDir
	cmd.Env = os.Environ()

	// if err := cmd.Start(); err != nil {
	// 	return err
	// }
	// if err := cmd.Process.Release(); err != nil {
	// 	return err
	// }
	// return nil

	out, err := cmd.Output()
	debugLogCmd(cmd, out, err)
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
	debugLogCmd(cmd, out, err)
	return err
}

func (r *runcRuntime) DeleteContainer() {
	panic("not implemented")
}

func (r *runcRuntime) ContainerState(id container.ID) (StateResp, error) {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"state",
		string(id),
	)
	out, err := cmd.Output()
	debugLogCmd(cmd, out, err)
	if err != nil {
		return StateResp{}, err
	}

	resp := StateResp{}
	return resp, json.Unmarshal(out, &resp)
}

func debugLogCmd(c *exec.Cmd, stdout []byte, err error) {
	logrus.WithFields(logrus.Fields{
		"stdout": string(stdout),
		"error":  err,
	}).Debugf("exec %s", strings.Join(c.Args, " "))
}
