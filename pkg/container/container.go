package container

import (
	"errors"
)

type Container struct {
	id      ID
	name    string
	rootfs  string
	command []string
	args    []string
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
	}, nil
}

func (c *Container) ID() ID {
	return c.id
}

func isValidName(name string) bool {
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 32
}
