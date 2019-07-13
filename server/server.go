package server

import (
	"golang.org/x/net/context"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/cri"
)

// Protobuf stuctures are completely hidden behind this abstraction.
type conmanServer struct {
	runtimeSrv cri.RuntimeService
}

type Conman interface {
	ConmanServer
	Serve(network, addr string) error
}

func New(runtimeSrv cri.RuntimeService) Conman {
	return &conmanServer{
		runtimeSrv: runtimeSrv,
	}
}

func (s *conmanServer) CreateContainer(
	ctx context.Context,
	req *CreateContainerRequest,
) (*CreateContainerResponse, error) {
	cont, err := s.runtimeSrv.CreateContainer(
		cri.ContainerOptions{
			Name:           req.Name,
			RootfsPath:     req.RootfsPath,
			RootfsReadonly: req.RootfsReadonly,
		},
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
	err := s.runtimeSrv.StartContainer(
		container.ID(req.ContainerId),
	)
	if err != nil {
		return nil, err
	}
	return &StartContainerResponse{}, nil
}
