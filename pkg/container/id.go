package container

import (
	"strings"

	"github.com/satori/go.uuid"
)

type ID string

func RandID() ID {
	return ID(strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
}
