# WHAT specifies the OS image to build
WHAT?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)

all: build
build:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null
	docker build -t ${DOCKER_USER}/ignite-${WHAT}:${VERSION} ${WHAT}

push:
	docker push ${DOCKER_USER}/ignite-${WHAT}:${VERSION}

build-all:
	make -C kernel
	$(MAKE) build WHAT=ubuntu
	$(MAKE) build WHAT=centos

push-all: build-all
	make -C kernel push
	$(MAKE) push WHAT=ubuntu
	$(MAKE) push WHAT=centos
