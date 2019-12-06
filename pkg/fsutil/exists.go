package fsutil

import (
	"os"

	"github.com/sirupsen/logrus"
)

func AssertExists(filename string) {
	ok, err := Exists(filename)
	if !ok || err != nil {
		logrus.WithError(err).Fatal("File is not reachable: " + filename)
	}
}

func EnsureExists(dir string) string {
	if err := os.MkdirAll(dir, 0755); err != nil {
		logrus.WithError(err).Fatal("File is not reachable: " + dir)
	}
	return dir
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
