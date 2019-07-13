package testutil

import (
	"log"

	"github.com/iximiuz/conman/pkg/container"
)

func NewContainer() *container.Container {
	id := container.RandID()
	name := "name_" + string(id[:8])
	c, err := container.New(id, name)
	if err != nil {
		log.Fatalf("Unexpected error during creation of test "+
			"container: %v\n id=%v name=%v\n", err, id, name)
	}
	return c
}
