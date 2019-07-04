ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY:
build_proto:
	docker run -it -v ${ROOT_DIR}:/opt/conman:rw grpc/go protoc -I/opt/conman/server conman.proto --go_out=plugins=grpc:/opt/conman/server

