package shimmy_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestMe(t *testing.T) {
	shimmy := "/home/vagrant/shimmy/target/debug/shimmy"
	cmd := exec.Command(shimmy)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	rw, wr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	cmd.ExtraFiles = append(cmd.ExtraFiles, wr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("_OCI_SYNCPIPE=%d", 3))
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	wr.Close()

	bytes, err := ioutil.ReadAll(rw)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bytes))
}
