package main

import (
	"github.com/iximiuz/conman/ctl/cmd"
	_ "github.com/iximiuz/conman/ctl/cmd/containers" // for init()
)

func main() {
	cmd.Execute()
}
