SHELL:=/bin/bash
# Set the command for running `docker`
# -- allows user to override for things like sudo usage or container images
DOCKER := docker
# Set the first containerd.sock that successfully stats -- fallback to the docker4mac default
CONTAINERD_SOCK := $(shell \
	$(DOCKER) run -i --rm \
		-v /run:/run:ro \
		-v /var/run:/var/run:ro \
		busybox:latest \
		ls 2>/dev/null \
		/run/containerd/containerd.sock \
		/run/docker/containerd/containerd.sock \
		/var/run/containerd/containerd.sock \
		/var/run/docker/containerd/containerd.sock \
		| head -n1 \
		|| echo \
			/var/run/docker/containerd/containerd.sock \
	)
# Set the command for running `ctr`
# Use root inside a container with the host containerd socket
# This is a form of privilege escalation that avoids interactive sudo during make
CTR := $(DOCKER) run -i --rm \
		-v $(CONTAINERD_SOCK):/run/containerd/containerd.sock \
		linuxkit/containerd:751de142273e1b5d2d247d2832d654ab92e907bc \
		ctr

# WHAT specifies the OS image to build
WHAT?=
IS_LATEST?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)
RELEASE?=latest
FULL_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:${TAG}
escaped_RELEASE:=$(subst +,-,$(RELEASE))
RELEASE_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:${escaped_RELEASE}
LATEST_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:latest

GOARCH?=amd64
GOARCH_LIST=amd64 arm64
# Set this flag to build a multiarch image
IS_MANIFEST_LIST?=0
TAG:=${escaped_RELEASE}$(if $(strip $(VERSION)),-${VERSION})
OP:=build

TMPDIR?=/tmp

ifeq ($(IS_MANIFEST_LIST),1)
BUILD_TMPDIR:=$(shell mktemp -d $(TMPDIR)/ignite-image-build.XXXXXXXXXX)
ARCH_TAG=-$(GOARCH)
endif

ifeq ($(GOARCH),amd64)
DOCKERARCH=amd64
endif
ifeq ($(GOARCH),arm64)
DOCKERARCH=arm64v8
QEMUARCH=aarch64
endif


all: build-all


build:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null

ifeq ($(IS_MANIFEST_LIST),0)
	sed "s|DOCKERARCH|$(DOCKERARCH)|g;/QEMUARCH/d" ${WHAT}/Dockerfile \
		| $(DOCKER) build \
			--build-arg RELEASE \
			--build-arg BINARY_REF \
			--build-arg GOARCH \
			--build-arg DIGEST \
			-f -\
			-t $(FULL_IMAGE_NAME) \
			${WHAT}
else
	# Register /usr/bin/qemu-ARCH-static as the handler for non-x86 binaries in the kernel
	$(DOCKER) run --rm --privileged multiarch/qemu-user-static:register --reset
	cp -a ${WHAT}/. $(BUILD_TMPDIR)/
	sed "s|QEMUARCH|$(QEMUARCH)|g;s|DOCKERARCH|$(DOCKERARCH)|g" ${WHAT}/Dockerfile > $(BUILD_TMPDIR)/Dockerfile

ifeq ($(GOARCH),amd64)
	sed "/COPY qemu/d" $(BUILD_TMPDIR)/Dockerfile > $(BUILD_TMPDIR)/Dockerfile.new
	mv $(BUILD_TMPDIR)/Dockerfile.new $(BUILD_TMPDIR)/Dockerfile
else
	$(MAKE) -C .. qemu
	cp ../bin/$(GOARCH)/qemu-$(QEMUARCH)-static $(BUILD_TMPDIR)
endif

	$(DOCKER) build \
		--build-arg RELEASE \
		--build-arg BINARY_REF \
		--build-arg GOARCH \
		--build-arg DIGEST \
		-t $(FULL_IMAGE_NAME)$(ARCH_TAG) \
		$(BUILD_TMPDIR)
	find $(BUILD_TMPDIR)/ -mindepth 1 -delete
endif

	$(DOCKER) tag $(FULL_IMAGE_NAME)$(ARCH_TAG) $(RELEASE_IMAGE_NAME)$(ARCH_TAG)
ifeq ($(IS_LATEST),true)
	$(DOCKER) tag $(FULL_IMAGE_NAME)$(ARCH_TAG) $(LATEST_IMAGE_NAME)$(ARCH_TAG)
endif


ctr-import:
	$(DOCKER) image save $(FULL_IMAGE_NAME)$(ARCH_TAG) \
		| $(CTR) -n firecracker image import -
	$(DOCKER) image save $(RELEASE_IMAGE_NAME)$(ARCH_TAG) \
		| $(CTR) -n firecracker image import -
ifeq ($(IS_LATEST),true)
	$(DOCKER) image save $(LATEST_IMAGE_NAME)$(ARCH_TAG) \
		| $(CTR) -n firecracker image import -
endif


push:
ifeq ($(IS_MANIFEST_LIST),0)
	$(DOCKER) push $(FULL_IMAGE_NAME)
	$(DOCKER) push $(RELEASE_IMAGE_NAME)
ifeq ($(IS_LATEST),true)
	$(DOCKER) push $(LATEST_IMAGE_NAME)
endif
else
ifeq ($(GOARCH),amd64)
	../hack/push-manifest-list.sh $(FULL_IMAGE_NAME) $(GOARCH_LIST)
	../hack/push-manifest-list.sh $(RELEASE_IMAGE_NAME) $(GOARCH_LIST)
ifeq ($(IS_LATEST),true)
	../hack/push-manifest-list.sh $(LATEST_IMAGE_NAME) $(GOARCH_LIST)
endif
endif
endif


ctr-import-all:
	$(MAKE) OP=ctr-import build-all


push-all: build-all
	$(MAKE) OP=push build-all


build-all: build-almalinux build-alpine build-amazon-kernel build-amazonlinux build-centos build-k3s build-kubeadm build-opensuse build-rockylinux build-ubuntu

build-almalinux:
	$(MAKE) ${OP} WHAT=almalinux      RELEASE=8                                        IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=amd64
	$(MAKE) ${OP} WHAT=almalinux      RELEASE=8                                        IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=arm64

build-alpine:
ifeq ($(OP),build)
	pushd alpine; \
		${MAKE} alpine.tar; \
	popd;
endif
	$(MAKE) ${OP} WHAT=alpine

build-amazon-kernel:
	$(MAKE) ${OP} WHAT=amazon-kernel

build-amazonlinux:
	$(MAKE) ${OP} WHAT=amazonlinux    RELEASE=2                                        IS_LATEST=true

build-centos:
	$(MAKE) ${OP} WHAT=centos         RELEASE=7
	$(MAKE) ${OP} WHAT=centos         RELEASE=8                                        IS_LATEST=true

build-k3s:
	$(MAKE) ${OP} WHAT=k3s            RELEASE=v1.20.4+k3s1                             IS_LATEST=true

build-kubeadm:
	$(MAKE) ${OP} WHAT=kubeadm        RELEASE=v1.18.3  BINARY_REF=release/stable-1.18  IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=arm64
	$(MAKE) ${OP} WHAT=kubeadm        RELEASE=v1.18.3  BINARY_REF=release/stable-1.18  IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=amd64

build-opensuse:
	$(MAKE) ${OP} WHAT=opensuse       RELEASE=leap                                     IS_LATEST=true
	$(MAKE) ${OP} WHAT=opensuse       RELEASE=tumbleweed

build-rockylinux:
	$(MAKE) ${OP} WHAT=rockylinux      RELEASE=8                                        IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=amd64  DIGEST=f0d7460b97156f6c8ea2ae73152bc11fe410d272387d60ddff36dfcea22ef689
	$(MAKE) ${OP} WHAT=rockylinux      RELEASE=8                                        IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=arm64  DIGEST=e98dd481db7973d701f8695df7be1b2393261cf273aa68a89d873b7b134588df

build-ubuntu:
	$(MAKE) ${OP} WHAT=ubuntu         RELEASE=16.04                                                    IS_MANIFEST_LIST=0
	$(MAKE) ${OP} WHAT=ubuntu         RELEASE=18.04                                                    IS_MANIFEST_LIST=1  GOARCH=arm64
	$(MAKE) ${OP} WHAT=ubuntu         RELEASE=18.04                                                    IS_MANIFEST_LIST=1  GOARCH=amd64
	$(MAKE) ${OP} WHAT=ubuntu         RELEASE=20.04                                    IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=arm64
	$(MAKE) ${OP} WHAT=ubuntu         RELEASE=20.04                                    IS_LATEST=true  IS_MANIFEST_LIST=1  GOARCH=amd64
