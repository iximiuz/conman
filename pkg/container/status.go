package container

import (
	"errors"
	"fmt"
	"math"
)

type Status uint32

const (
	Initial Status = 0
	Created Status = 10
	Running Status = 20
	Stopped Status = 30
	Unknown Status = math.MaxUint32
)

func StatusFromString(s string) (Status, error) {
	switch s {
	case "created":
		return Created, nil
	case "running":
		return Running, nil
	case "stopped":
		return Stopped, nil
	}
	return Unknown, errors.New(fmt.Sprintf("Unknown status %s", s))
}

func (s Status) String() string {
	switch s {
	case Initial:
		return "initial"
	case Created:
		return "created"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	}
	panic("unreachable")
}
