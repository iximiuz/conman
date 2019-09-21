package container_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/iximiuz/conman/pkg/container"
)

func TestMarshalUnmarshalJSON(t *testing.T) {
	bytes1 := []byte(`{"id":"1","name":"cont1","status":0,"createdAt":"2019-09-21T14:35:29Z","command":["/bin/sh"],"rootfs":"/path/to/bundle"}`)

	cont := &container.Container{}
	if err := json.Unmarshal(bytes1, cont); err != nil {
		t.Fatal(err)
	}

	if cont.ID() != "1" {
		t.Fatal("Unexpected ID")
	}
	if cont.Name() != "cont1" {
		t.Fatal("Unexpected name")
	}
	if cont.Status() != container.Initial {
		t.Fatal("Unexpected status")
	}
	if cont.CreatedAt() != "2019-09-21T14:35:29Z" {
		t.Fatal("Unexpected createdAt")
	}

	bytes2, err := json.Marshal(cont)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes1, bytes2) {
		t.Fatal("Marshal(Unmarshal(b)) != Unmarshal(Marshal(c))")
	}
}
