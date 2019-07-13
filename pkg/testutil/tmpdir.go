package testutil

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func TempDir(t *testing.T, suffix ...string) string {
	dir, err := ioutil.TempDir("", t.Name()+strings.Join(suffix, "_"))
	if err != nil {
		log.Fatal("Cannot create tmp dir")
	}
	return dir
}
