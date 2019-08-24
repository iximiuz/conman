package testutil

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
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

func DataDir(name string) string {
	p, err := filepath.Abs(path.Join("../../test/data", name))
	if err != nil {
		log.Fatal("Cannot create tmp dir")
	}
	return p
}
