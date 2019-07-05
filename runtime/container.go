package runtime

import (
	"errors"
)

type ContainerID string

type Container struct {
	id      ContainerID
	name    string
	rootfs  string
	command []string
	args    []string
}

func NewContainer(
	id ContainerID,
	name string,
) (*Container, error) {
	if !isValidContainerName(name) {
		return nil, errors.New("Invalid container name")
	}

	return &Container{
		id:   id,
		name: name,
	}, nil
}

func (c *Container) ID() ContainerID {
	return c.id
}

func isValidContainerName(name string) bool {
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 32
}
