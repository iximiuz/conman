package testutil

import (
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func ReadPidfile(filename string) (pid int, err error) {
	if bytes, err := ioutil.ReadFile(filename); err == nil {
		pid, err = strconv.Atoi(string(bytes))
	}
	return
}

func FindProcessByPidfile(filename string) (*os.Process, error) {
	pid, err := ReadPidfile(filename)
	if err != nil {
		return nil, err
	}
	return os.FindProcess(pid)
}

func EnsureProcessIsAlive(proc *os.Process) error {
	return proc.Signal(syscall.Signal(0))
}

func EnsureProcessHasTerminated(proc *os.Process, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			if err.Error() == "os: process already finished" {
				return nil
			}
			return err
		}
		time.Sleep(100 * time.Microsecond)
	}

	return errors.Errorf("Process %v is still alive after %v.", proc.Pid, timeout)
}
