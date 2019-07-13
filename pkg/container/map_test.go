package container_test

import (
	"testing"

	"github.com/iximiuz/conman/pkg/container"
	"github.com/iximiuz/conman/pkg/rollback"
	"github.com/iximiuz/conman/pkg/testutil"
)

func TestAdd(t *testing.T) {
	if err := container.NewMap().Add(testutil.NewContainer(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestAddWithRollback(t *testing.T) {
	m := container.NewMap()
	c := testutil.NewContainer()
	rb := rollback.New()
	if err := m.Add(c, rb); err != nil {
		t.Fatal(err)
	}

	if m.Get(c.ID()) == nil {
		t.Fatal("Container not found")
	}

	rb.Execute()

	if m.Get(c.ID()) != nil {
		t.Fatal("Container has not been deleted by rollback")
	}
}

func TestGet(t *testing.T) {
	m := container.NewMap()
	c := testutil.NewContainer()

	if err := m.Add(c, nil); err != nil {
		t.Fatal("Unexpected", err)
	}
	if m.Get(c.ID()) == nil {
		t.Fatal("Get() returned nil")
	}
}

func TestGetByName(t *testing.T) {
	m := container.NewMap()
	c := testutil.NewContainer()

	if err := m.Add(c, nil); err != nil {
		t.Fatal("Unexpected", err)
	}
	if m.GetByName(c.Name()) == nil {
		t.Fatal("GetByName() returned nil")
	}
}
