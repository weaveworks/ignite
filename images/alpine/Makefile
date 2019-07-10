TEMP_DIR:=$(shell mktemp -d)
VERSION?=$(shell git describe HEAD --tags)
DOCKER_USER?=weaveworks

all: build

alpine.ext4:
	curl -sSL https://s3.amazonaws.com/spec.ccfc.min/img/hello/fsfiles/hello-rootfs.ext4 > alpine.ext4

alpine.tar: alpine.ext4
	sudo mount -o loop $(shell pwd)/alpine.ext4 ${TEMP_DIR}
	cd ${TEMP_DIR} && sudo tar -cf $(shell pwd)/alpine.tar .
	sudo umount $(shell pwd)/alpine.ext4
#  alpine.tar
build:
	docker build -t ${DOCKER_USER}/ignite-alpine:${VERSION} .

push: build
	docker push ${DOCKER_USER}/ignite-alpine:${VERSION}
	docker tag ${DOCKER_USER}/ignite-alpine:${VERSION} ${DOCKER_USER}/ignite-alpine:latest
