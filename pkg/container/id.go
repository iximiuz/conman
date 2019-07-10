package container

import (
	"github.com/satori/go.uuid"
)

type ID string

func RandID() ID {
	return ID(uuid.NewV4().String())
}
