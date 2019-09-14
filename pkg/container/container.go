package container

import (
	"errors"
	"time"
)

type Container struct {
	id      ID
	name    string
	state   state
	rootfs  string
	command []string
	args    []string
}

type state struct {
	pid        int
	status     Status
	createdAt  string
	startedAt  string
	finishedAt string
}

func New(
	id ID,
	name string,
) (*Container, error) {
	if !isValidName(name) {
		return nil, errors.New("Invalid container name")
	}

	return &Container{
		id:   id,
		name: name,
		state: state{
			createdAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (c *Container) ID() ID {
	return c.id
}

func (c *Container) Name() string {
	return c.name
}

func (c *Container) CreatedAt() string {
	return c.state.createdAt
}

func (c *Container) CreatedAtNano() int64 {
	return unixNanoTime(c.CreatedAt())
}

func (c *Container) StartedAt() string {
	return c.state.startedAt
}

func (c *Container) StartedAtNano() int64 {
	if c.state.startedAt == "" {
		return 0
	}
	return unixNanoTime(c.StartedAt())
}

func (c *Container) FinishedAt() string {
	return c.state.finishedAt
}

func (c *Container) FinishedAtNano() int64 {
	if c.state.finishedAt == "" {
		return 0
	}
	return unixNanoTime(c.FinishedAt())
}

func (c *Container) Status() Status {
	return c.state.status
}

func (c *Container) SetStatus(s Status) {
	c.state.status = s
}

func isValidName(name string) bool {
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && (c != '_') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 32
}

func unixNanoTime(s string) int64 {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.UnixNano()
}
