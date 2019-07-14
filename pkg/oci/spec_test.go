package oci

import (
	"testing"
)

func TestNewSpec(t *testing.T) {
	spec, err := NewSpec(SpecOptions{})
	if err != nil {
		t.Fatal("NewSpec() failed", err)
	}
	t.Log(len(spec))
}
