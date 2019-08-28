package oci

import (
	"errors"
	"os"
	"syscall"
)

var signals = map[os.Signal]string{
	syscall.SIGKILL: "KILL",
	syscall.SIGTERM: "TERM",
}

func sigStr(sig os.Signal) (string, error) {
	if str, ok := signals[sig]; ok {
		return str, nil
	}
	return "", errors.New("Unknown signal")
}
