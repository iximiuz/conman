ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY:
bin/conmand:
	go build -o $@ ${ROOT_DIR}/main.go

.PHONY:
bin/conmanctl:
	go build -o $@ ${ROOT_DIR}/ctl/main.go

.PHONY:
build_proto:
	docker run -it -v ${ROOT_DIR}:/opt/conman:rw grpc/go protoc -I/opt/conman/server conman.proto --go_out=plugins=grpc:/opt/conman/server

