package container

import (
	"errors"
	"sync"

	"github.com/iximiuz/conman/pkg/rollback"
)

type Map struct {
	sync.RWMutex
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
	m.Lock()
	defer m.Unlock()

	if _, ok := m.byid[c.ID()]; ok {
		return errors.New("Duplicate container ID")
	}
	if _, ok := m.byname[c.Name()]; ok {
		return errors.New("Duplicate container name")
	}

	m.byid[c.ID()] = c
	m.byname[c.Name()] = c

	if rb != nil {
		rb.Add(func() { m.Del(c.ID()) })
	}
	return nil
}

func (m *Map) Get(id ID) *Container {
	c, _ := m.byid[id]
	return c
}

func (m *Map) GetByName(name string) *Container {
	c, _ := m.byname[name]
	return c
}

func (m *Map) Del(id ID) bool {
	c, ok := m.byid[id]
	if ok {
		delete(m.byid, id)
		delete(m.byname, c.Name())
	}
	return ok
}
