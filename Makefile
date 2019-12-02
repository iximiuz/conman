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

.PHONY:
testfunctional:
	bash test/conman/runner.sh

test/data/rootfs_alpine:
	$(eval CID=$(shell docker create -l com.iximiuz-project=${REPO} alpine))
	mkdir -p ${ROOT_DIR}/test/data/rootfs_alpine/
	docker export ${CID} | tar -C ${ROOT_DIR}/test/data/rootfs_alpine/ -xvf -
	docker rm ${CID}

.PHONY:
build_proto:
	docker run --rm -v ${ROOT_DIR}:/opt/conman:rw grpc/go protoc -I/opt/conman/server conman.proto --go_out=plugins=grpc:/opt/conman/server

.PHONY:
clean: clean-docker-procs clean-lib-root clean-test-runs

.PHONY:
clean-docker-procs:
	@echo "[Remove Docker Processes]"
	@if [ "`docker ps -qa -f=label=com.iximiuz-project=${REPO}`" != '' ]; then\
		docker rm `docker ps -qa -f=label=com.iximiuz-project=${REPO}`;\
	else\
		echo "<noop>";\
	fi

.PHONY:
clean-lib-root:
	@echo "[Remove conman lib directory]"
	rm -rf /var/lib/conman

.PHONY:
clean-test-runs:
	@echo "[Clean test runs]"
	rm -f "${ROOT_DIR}/test/conman/conmand.log"
	find /tmp -name conman-test-run.* -type d | xargs rm -rfv

