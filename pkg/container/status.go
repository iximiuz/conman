package container

import (
	"errors"
	"fmt"
	"math"
)

type Status uint32

const (
	StatusNew     = 0
	StatusCreated = 10
	StatusRunning = 20
	StatusUnknown = math.MaxUint32
)

func StatusFromString(s string) (Status, error) {
	switch s {
	case "created":
		return StatusCreated, nil
	}
	return StatusUnknown, errors.New(fmt.Sprintf("Unknown status %s", s))
}

func (s Status) String() string {
	switch s {
	case StatusNew:
		return "new"
	case StatusCreated:
		return "created"
	case StatusRunning:
		return "run"
	}
	panic("unreachable")
}
