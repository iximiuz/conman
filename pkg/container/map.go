package container

import "github.com/iximiuz/conman/pkg/rollback"

// TODO: add mutex!
type Map struct {
	byid   map[ID]*Container
	byname map[string]*Container
}

func NewMap() *Map {
	return &Map{
		byid:   make(map[ID]*Container),
		byname: make(map[string]*Container),
	}
}

func (m *Map) Add(c *Container, rb *rollback.Rollback) error {
	if rb != nil {
		rb.Add(func() { m.Del(c.ID()) })
	}
	return nil
}

func (m *Map) Get(id ID) *Container {
	return nil
}

func (m *Map) GetByName(name string) *Container {
	return nil
}

func (m *Map) Del(id ID) bool {
	return false
}
