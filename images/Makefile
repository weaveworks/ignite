# WHAT specifies the OS image to build
WHAT?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)
RELEASE?=latest
TAG:=${RELEASE}$(if $(strip $(VERSION)),-${VERSION})

all: build
build:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null
	docker build --build-arg RELEASE -t ${DOCKER_USER}/ignite-${WHAT}:${TAG} ${WHAT}

push:
	docker push ${DOCKER_USER}/ignite-${WHAT}:${TAG}
	docker tag ${DOCKER_USER}/ignite-${WHAT}:${TAG} ${DOCKER_USER}/ignite-${WHAT}:latest
	docker push ${DOCKER_USER}/ignite-${WHAT}:latest

build-all:
	# Temporarily comment out the kernel build for the RC
	#	make -C kernel
	$(MAKE) build WHAT=amazon-kernel
	$(MAKE) build WHAT=amazonlinux
	$(MAKE) build WHAT=alpine
	$(MAKE) build WHAT=ubuntu
	$(MAKE) build WHAT=centos
	$(MAKE) build WHAT=kubeadm

push-all: build-all
	# Temporarily comment out the kernel build for the RC
	#	make -C kernel push
	$(MAKE) push WHAT=amazon-kernel
	$(MAKE) push WHAT=amazonlinux
	$(MAKE) push WHAT=alpine
	$(MAKE) push WHAT=ubuntu
	$(MAKE) push WHAT=centos
	$(MAKE) push WHAT=kubeadm
