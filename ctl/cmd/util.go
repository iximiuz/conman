package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/iximiuz/conman/server"
)

func Print(v interface{}) {
	s, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(string(s))
}

func Connect() (server.ConmanClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("unix://"+OptHost, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	return server.NewConmanClient(conn), conn
}
