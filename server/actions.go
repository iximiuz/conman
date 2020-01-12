package server

import (
	"time"

	"golang.org/x/net/context"
	criapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/cri"
)

func (s *conmanServer) Version(
	ctx context.Context,
	req *VersionRequest,
) (resp *VersionResponse, err error) {
	return &VersionResponse{
		Version:     "0.0.1",
		RuntimeName: "runc",
	}, nil
}

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

func (s *conmanServer) StopContainer(
	ctx context.Context,
	req *StopContainerRequest,
) (resp *StopContainerResponse, err error) {
	traceRequest("StopContainer", req)
	defer func() { traceResponse("StopContainer", resp, err) }()

	err = s.runtimeSrv.StopContainer(
		container.ID(req.ContainerId),
		time.Duration(req.Timeout)*time.Second,
	)
	if err == nil {
		resp = &StopContainerResponse{}
	}
	return
}

func (s *conmanServer) RemoveContainer(
	ctx context.Context,
	req *RemoveContainerRequest,
) (resp *RemoveContainerResponse, err error) {
	traceRequest("RemoveContainer", req)
	defer func() { traceResponse("RemoveContainer", resp, err) }()

	err = s.runtimeSrv.RemoveContainer(
		container.ID(req.ContainerId),
	)
	if err == nil {
		resp = &RemoveContainerResponse{}
	}
	return
}

func (s *conmanServer) ListContainers(
	ctx context.Context,
	req *ListContainersRequest,
) (resp *ListContainersResponse, err error) {
	traceRequest("ListContainers", req)
	defer func() { traceResponse("ListContainers", resp, err) }()

	cs, err := s.runtimeSrv.ListContainers()
	if err != nil {
		return nil, err
	}

	return &ListContainersResponse{
		Containers: toPbContainers(cs),
	}, nil
}

func (s *conmanServer) ContainerStatus(
	ctx context.Context,
	req *ContainerStatusRequest,
) (resp *ContainerStatusResponse, err error) {
	traceRequest("ContainerStatus", req)
	defer func() { traceResponse("ContainerStatus", resp, err) }()

	cont, err := s.runtimeSrv.GetContainer(
		container.ID(req.ContainerId),
	)
	if err != nil {
		return nil, err
	}

	return &ContainerStatusResponse{
		Status: &ContainerStatus{
			ContainerId:   string(cont.ID()),
			ContainerName: string(cont.Name()),
			State:         toPbContainerState(cont.Status()),
			CreatedAt:     cont.CreatedAtNano(),
			StartedAt:     cont.StartedAtNano(),
			FinishedAt:    cont.FinishedAtNano(),
			ExitCode:      cont.ExitCode(),
			LogPath:       cont.LogPath(),
		},
	}, nil
}

func (s *conmanServer) Attach(
	ctx context.Context,
	req *AttachRequest,
) (resp *AttachResponse, err error) {
	traceRequest("Attach", req)
	defer func() { traceResponse("Attach", resp, err) }()

	r, err := s.streamingSrv.GetAttach(&criapi.AttachRequest{
		ContainerId: req.ContainerId,
		Tty:         req.Tty,
		Stdin:       req.Stdin,
		Stdout:      req.Stdout,
		Stderr:      req.Stderr,
	})
	if err != nil {
		return nil, err
	}
	return &AttachResponse{Url: r.Url}, err
}

func toPbContainerState(s container.Status) ContainerState {
	switch s {
	case container.Created:
		return ContainerState_CREATED
	case container.Running:
		return ContainerState_RUNNING
	case container.Stopped:
		return ContainerState_EXITED
	}
	return ContainerState_UNKNOWN
}

func toPbContainers(cs []*container.Container) (rv []*Container) {
	for _, c := range cs {
		rv = append(rv, &Container{
			Id:        string(c.ID()),
			Name:      string(c.Name()),
			CreatedAt: c.CreatedAtNano(),
			State:     toPbContainerState(c.Status()),
		})
	}
	return
}
