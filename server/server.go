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
	ctx context.Context,
	req *CreateContainerRequest,
) (*CreateContainerResponse, error) {
	cont, err := s.rt.CreateContainer(
		req.Name,
	)
	if err != nil {
		return nil, err
	}
	return &CreateContainerResponse{
		ContainerId: string(cont.ID()),
	}, nil
}

func (s *conmanServer) StartContainer(
	ctx context.Context,
	req *StartContainerRequest,
) (*StartContainerResponse, error) {
	err := s.rt.StartContainer(runtime.ContainerID(req.ContainerId))
	if err != nil {
		return nil, err
	}
	return &StartContainerResponse{}, nil
}
