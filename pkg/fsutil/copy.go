package fsutil

import (
	"os/exec"
)

func CopyDir(src, dst string) error {
	cmd := exec.Command("cp", "-a", src, dst)
	return cmd.Run()
}
