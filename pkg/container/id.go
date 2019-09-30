package container

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/satori/go.uuid"
)

type ID string

var badIdFormatErr = errors.New("Bad container ID format")

func RandID() ID {
	return ID(strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
}

func ParseID(id string) (ID, error) {
	if len(id) != 32 {
		return ID(""), badIdFormatErr
	}
	if _, err := hex.DecodeString(id); err != nil {
		return ID(""), badIdFormatErr
	}
	return ID(id), nil
}
