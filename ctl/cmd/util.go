package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func Print(v interface{}) {
	s, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(string(s))
}
