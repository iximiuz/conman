package cri

import (
	"io"
	"net"
	"path"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/iximiuz/conman/pkg/container"
)

func (rs *runtimeService) Attach(
	containerID string,
	stdin io.Reader,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
	tty bool,
	resize <-chan remotecommand.TerminalSize,
) error {
	go func() {
		for i := 0; i < 10; i++ {
			stdout.Write([]byte("Hi there!"))
			time.Sleep(time.Second)
		}
	}()

	cont, err := rs.GetContainer(container.ID(containerID))
	if err != nil {
		return err
	}
	// TODO: check cont.Status() in [CREATED,RUNNING]

	hcont, err := rs.cstore.GetContainer(cont.ID())
	if err != nil {
		return err
	}

	attachSocketPath := path.Join(hcont.BundleDir(), "attach")
	conn, err := net.DialUnix(
		"unixpacket",
		nil,
		&net.UnixAddr{Name: attachSocketPath, Net: "unixpacket"},
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	time.Sleep(12 * time.Second)
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
