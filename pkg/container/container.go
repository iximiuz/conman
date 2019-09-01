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
	pid       int
	status    Status
	createdAt string
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
