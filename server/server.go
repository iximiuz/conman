package server

import (
	"golang.org/x/net/context"

	"github.com/iximiuz/conman/runtime"
)

// Protobuf stuctures are completely hidden behind this abstraction.
type conmanServer struct {
	rt runtime.Runtime
}

type Conman interface {
	ConmanServer
	Serve(network, addr string) error
}

func New(rt runtime.Runtime) Conman {
	return &conmanServer{
		rt: rt,
	}
}

func (s *conmanServer) CreateContainer(
	_ctx context.Context,
	_req *CreateContainerRequest,
) (*CreateContainerResponse, error) {
	_, err := s.rt.CreateContainer("")
	if err != nil {
		return nil, err
	}
	return &CreateContainerResponse{}, nil
}

func (s *conmanServer) StartContainer(
	_ctx context.Context,
	_req *StartContainerRequest,
) (*StartContainerResponse, error) {
	_, err := s.rt.StartContainer("")
	if err != nil {
		return nil, err
	}
	return &StartContainerResponse{}, nil
}
