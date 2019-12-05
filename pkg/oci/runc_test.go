package oci_test

import (
	"testing"

	"github.com/iximiuz/conman/config"
)

var cfg *config.Config

func init() {
	cfg = config.TestConfigFromFlags()
}

func Test_NonInteractive_SimpleRun(t *testing.T) {

}
