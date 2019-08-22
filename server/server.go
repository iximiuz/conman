package server

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/iximiuz/conman/pkg/cri"
)

type Server interface {
	ConmanServer // pb server interface
	Serve(network, addr string) error
}

// Protobuf stuctures are completely hidden behind this abstraction.
type conmanServer struct {
	runtimeSrv cri.RuntimeService
}

func New(runtimeSrv cri.RuntimeService) Server {
	return &conmanServer{
		runtimeSrv: runtimeSrv,
	}
}

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

func traceRequest(name string, req interface{}) {
	logrus.WithFields(logrus.Fields{
		"body": req,
	}).Trace("Request [" + name + "]")
}

func traceResponse(name string, resp interface{}, err error) {
	logrus.WithFields(logrus.Fields{
		"body": resp,
	}).WithError(err).Trace("Response [" + name + "]")
}
