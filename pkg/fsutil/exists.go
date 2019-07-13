package fsutil

import "os"

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
