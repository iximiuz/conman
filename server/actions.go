package server

import (
	"golang.org/x/net/context"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/cri"
)

func (s *conmanServer) CreateContainer(
	ctx context.Context,
	req *CreateContainerRequest,
) (resp *CreateContainerResponse, err error) {
	traceRequest("CreateContainer", req)
	defer func() { traceResponse("CreateContainer", resp, err) }()

	cont, err := s.runtimeSrv.CreateContainer(
		cri.ContainerOptions{
			Name:           req.Name,
			Command:        req.Command,
			Args:           req.Args,
			RootfsPath:     req.RootfsPath,
			RootfsReadonly: req.RootfsReadonly,
		},
	)
	if err == nil {
		resp = &CreateContainerResponse{
			ContainerId: string(cont.ID()),
		}
	}
	return
}

func (s *conmanServer) StartContainer(
	ctx context.Context,
	req *StartContainerRequest,
) (resp *StartContainerResponse, err error) {
	traceRequest("StartContainer", req)
	defer func() { traceResponse("StartContainer", resp, err) }()

	err = s.runtimeSrv.StartContainer(
		container.ID(req.ContainerId),
	)
	if err == nil {
		resp = &StartContainerResponse{}
	}
	return
}
