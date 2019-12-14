package fsutil

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func AssertExists(filename string) string {
	ok, err := Exists(filename)
	if !ok || err != nil {
		logrus.WithError(err).Fatal("File is not reachable: " + filename)
	}
	return filename
}

func EnsureExists(dirs ...string) string {
	target := path.Join(dirs...)
	if err := os.MkdirAll(target, 0755); err != nil {
		logrus.WithError(err).Fatal("Directory is not reachable: " + target)
	}
	return target
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
