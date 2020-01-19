package cri

import (
	"bytes"
	"io"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/iximiuz/conman/pkg/container"
)

// Same as in shimmy
const BufSize = 32 * 1024
const PipeTypeStdout = 1
const PipeTypeStderr = 2

func (rs *runtimeService) Attach(
	containerID string,
	stdin io.Reader,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
	_tty bool,
	_resize <-chan remotecommand.TerminalSize,
) error {
	if stdin == nil && stdout == nil && stderr == nil {
		return errors.New("at least one of the std streams must be open")
	}

	cont, err := rs.GetContainer(container.ID(containerID))
	if err != nil {
		return err
	}
	if cont.Status() != container.Created && cont.Status() != container.Running {
		return errors.Errorf("cannot connect to container in status %v", cont.Status())
	}

	conn, err := net.DialUnix(
		"unix",
		nil,
		&net.UnixAddr{Name: rs.containerAttachFile(cont.ID()), Net: "unix"},
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	doneOut := make(chan error)
	if stdout != nil || stderr != nil {
		go func() {
			doneOut <- forwardOutStreams(conn, stdout, stderr)
		}()
	}

	doneIn := make(chan error)
	if stdin != nil {
		go func() {
			_, err1 := io.Copy(conn, stdin)
			err2 := conn.CloseWrite()
			if err2 != nil {
				doneIn <- err2
			}
			doneIn <- err1
		}()
	}

	select {
	case err := <-doneIn:
		if stdout != nil || stderr != nil {
			return <-doneOut
		}
		return err
	case err := <-doneOut:
		return err
	}
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

func forwardOutStreams(conn io.Reader, stdout, stderr io.Writer) error {
	buf := make([]byte, BufSize+1)

	for {
		nread, err := conn.Read(buf)
		if nread > 0 {
			var dst io.Writer
			switch buf[0] {
			case PipeTypeStdout:
				dst = stdout
			case PipeTypeStderr:
				dst = stderr
			default:
				logrus.Debugf("unexpected attach pipe type %+d", buf[0])
			}

			if dst != nil {
				src := bytes.NewReader(buf[1:nread])
				if _, err := io.Copy(dst, src); err != nil {
					return err
				}
			}
		}
		if err == io.EOF || nread == 0 {
			return nil
		}
		if err != nil {
			return err
		}
	}

	// Unreachable!
	return nil
}
