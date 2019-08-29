package container_test

import (
	"testing"

	. "github.com/iximiuz/conman/pkg/container"
)

func TestStatusFromString(t *testing.T) {
	assertStatus(t, StatusCreated)(StatusFromString("created"))
	assertStatus(t, StatusUnknown)(StatusFromString("foobar"))
}

func TestStatusToString(t *testing.T) {
	assertString(t, StatusNew, "new")
	assertString(t, StatusCreated, "created")
}

func assertString(t *testing.T, s Status, expected string) {
	actual := s.String()
	if expected != actual {
		t.Fatalf("Status string mismatch: expected=%v actual=%v",
			expected, actual)
	}
}

func assertStatus(
	t *testing.T,
	expected Status,
) func(actual Status, err error) {
	return func(actual Status, err error) {
		if expected == actual {
			if expected != StatusUnknown && err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("Status mismatch: expected=%v actual=%v err=%v",
				expected, actual, err)
		}
	}
}
