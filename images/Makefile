# WHAT specifies the OS image to build
WHAT?=
IS_LATEST?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)
RELEASE?=latest
TAG:=${RELEASE}$(if $(strip $(VERSION)),-${VERSION})
OP:=build

all: build
build:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null
	docker build --build-arg RELEASE -t ${DOCKER_USER}/ignite-${WHAT}:${TAG} ${WHAT}

push:
	docker push ${DOCKER_USER}/ignite-${WHAT}:${TAG}
	docker tag ${DOCKER_USER}/ignite-${WHAT}:${TAG} ${DOCKER_USER}/ignite-${WHAT}:${RELEASE}
	docker push ${DOCKER_USER}/ignite-${WHAT}:${RELEASE}
ifeq ($(IS_LATEST),true)
	docker tag ${DOCKER_USER}/ignite-${WHAT}:${TAG} ${DOCKER_USER}/ignite-${WHAT}:latest
	docker push ${DOCKER_USER}/ignite-${WHAT}:latest
endif

push-all: build-all
	$(MAKE) OP=push build-all

build-all:
	$(MAKE) ${OP} WHAT=amazon-kernel
	$(MAKE) ${OP} WHAT=amazonlinux		RELEASE=2		IS_LATEST=true
	$(MAKE) ${OP} WHAT=alpine
	$(MAKE) ${OP} WHAT=ubuntu			RELEASE=16.04
	$(MAKE) ${OP} WHAT=ubuntu			RELEASE=18.04	IS_LATEST=true
	$(MAKE) ${OP} WHAT=ubuntu			RELEASE=19.04
	$(MAKE) ${OP} WHAT=centos			RELEASE=7		IS_LATEST=true
	$(MAKE) ${OP} WHAT=kubeadm