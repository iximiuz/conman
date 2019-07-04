package server

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

func (s *conmanServer) Serve(network, addr string) error {
	lis, err := listen(network, addr)
	if err != nil {
		return err
	}

	gsrv := grpc.NewServer()
	RegisterConmanServer(gsrv, s)
	return gsrv.Serve(lis)
}

func listen(network, addr string) (net.Listener, error) {
	if network != "unix" {
		return nil, errors.New("Only UNIX sockets supported")
	}
	if err := os.MkdirAll(filepath.Dir(addr), 0755); err != nil {
		return nil, err
	}
	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return net.Listen("unix", addr)
}
