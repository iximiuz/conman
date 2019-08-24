REPO=conman
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY:
bin/conmand:
	go build -o $@ ${ROOT_DIR}/main.go

.PHONY:
bin/conmanctl:
	go build -o $@ ${ROOT_DIR}/ctl/main.go

.PHONY:
testunit:
	go test ./...

test/data/rootfs_alpine:
	$(eval CID=$(shell docker create -l com.iximiuz-project=${REPO} alpine))
	mkdir -p ${ROOT_DIR}/test/data/rootfs_alpine/
	docker export ${CID} | tar -C ${ROOT_DIR}/test/data/rootfs_alpine/ -xvf -
	docker rm ${CID}

.PHONY:
build_proto:
	docker run -it -v ${ROOT_DIR}:/opt/conman:rw grpc/go protoc -I/opt/conman/server conman.proto --go_out=plugins=grpc:/opt/conman/server

.PHONY:
clean: clean-docker-procs

.PHONY:
clean-docker-procs:
	@echo "[Remove Docker Processes]"
	@if [ "`docker ps -qa -f=label=com.iximiuz-project=${REPO}`" != '' ]; then\
		docker rm `docker ps -qa -f=label=com.iximiuz-project=${REPO}`;\
	else\
		echo "<noop>";\
	fi

