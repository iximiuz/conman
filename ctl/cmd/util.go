package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/iximiuz/conman/server"
)

func Print(v interface{}) {
	fmt.Println(toString(v))
}

func Connect() (server.ConmanClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("unix://"+OptHost, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	return server.NewConmanClient(conn), conn
}

func toString(v interface{}) string {
	switch i := v.(type) {
	case proto.Message:
		s, err := (&jsonpb.Marshaler{EmitDefaults: true}).MarshalToString(i)
		if err != nil {
			logrus.Fatal(err)
		}
		return s
	default:
		s, err := json.Marshal(i)
		if err != nil {
			logrus.Fatal(err)
		}
		return string(s)
	}
}
