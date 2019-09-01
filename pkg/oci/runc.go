package oci

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	// "syscall"

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
	attrs := os.ProcAttr{
		Dir:   bundleDir,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		// Sys:   &syscall.SysProcAttr{Noctty: true}, TODO: ...
	}
	args := []string{r.exePath, "--root", r.rootPath, "start", string(id)}
	proc, err := os.StartProcess(args[0], args, &attrs)
	fmt.Printf("%+v\n%+v\n", proc, err)
	if err != nil {
		return err
	}
	return proc.Release()
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
