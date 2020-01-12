package cri

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/remotecommand"
)

func (rs *runtimeService) Attach(
	containerID string,
	in io.Reader,
	out io.WriteCloser,
	err io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
) error {
	go func() {
		for i := 0; i < 10; i++ {
			out.Write([]byte("Hi there!"))
			time.Sleep(time.Second)
		}
	}()
	return nil
}

func (rs *runtimeService) Exec(
	containerID string,
	cmd []string,
	in io.Reader,
	out io.WriteCloser,
	err io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
) error {
	return errors.New("Not implemented")
}

func (rs *runtimeService) PortForward(
	podSandboxID string,
	port int32,
	stream io.ReadWriteCloser,
) error {
	return errors.New("Not implemented")
}
