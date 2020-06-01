# WHAT specifies the OS image to build
WHAT?=
IS_LATEST?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)
RELEASE?=latest
FULL_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:${TAG}
RELEASE_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:${RELEASE}
LATEST_IMAGE_NAME=${DOCKER_USER}/ignite-${WHAT}:latest

GOARCH?=amd64
GOARCH_LIST=amd64 arm64
# Set this flag to build a multiarch image
IS_MANIFEST_LIST?=0
TAG:=${RELEASE}$(if $(strip $(VERSION)),-${VERSION})
OP:=build

ifeq ($(IS_MANIFEST_LIST),1)
TEMP_DIR:=$(shell mktemp -d)
ARCH_TAG=-$(GOARCH)
endif

ifeq ($(GOARCH),amd64)
DOCKERARCH=amd64
endif
ifeq ($(GOARCH),arm64)
DOCKERARCH=arm64v8
QEMUARCH=aarch64
endif

all: build
build:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null

ifeq ($(IS_MANIFEST_LIST),0)
	sed "s|DOCKERARCH|$(DOCKERARCH)|g;/QEMUARCH/d" ${WHAT}/Dockerfile | docker build --build-arg RELEASE -f - -t $(FULL_IMAGE_NAME) ${WHAT}
else
	# Register /usr/bin/qemu-ARCH-static as the handler for non-x86 binaries in the kernel
	docker run --rm --privileged multiarch/qemu-user-static:register --reset
	sed "s|QEMUARCH|$(QEMUARCH)|g;s|DOCKERARCH|$(DOCKERARCH)|g" ${WHAT}/Dockerfile > $(TEMP_DIR)/Dockerfile
ifeq ($(GOARCH),amd64)
	sed "/COPY qemu/d" -i $(TEMP_DIR)/Dockerfile
else
	$(MAKE) -C .. qemu
	cp ../bin/$(GOARCH)/qemu-$(QEMUARCH)-static $(TEMP_DIR)
endif
	docker build --build-arg RELEASE -t $(FULL_IMAGE_NAME)$(ARCH_TAG) $(TEMP_DIR)
endif

	docker tag $(FULL_IMAGE_NAME)$(ARCH_TAG) $(RELEASE_IMAGE_NAME)$(ARCH_TAG)
ifeq ($(IS_LATEST),true)
	docker tag $(FULL_IMAGE_NAME)$(ARCH_TAG) $(LATEST_IMAGE_NAME)$(ARCH_TAG)
endif

build-for-firekube:
ifeq ($(WHAT),)
	$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null
	(cd ${WHAT} && bash -c ./pull_images.sh)
	docker build --build-arg RELEASE -t ${DOCKER_USER}/ignite-centos:${RELEASE}-${SUB_TAG} ${WHAT}

build-image-firekube:
	$(MAKE) build-for-firekube WHAT=centos-firekube RELEASE=7 SUB_TAG=wksctl-v0.8.1

push:
ifeq ($(IS_MANIFEST_LIST),0)
	docker push $(FULL_IMAGE_NAME)
	docker push $(RELEASE_IMAGE_NAME)
ifeq ($(IS_LATEST),true)
	docker push $(LATEST_IMAGE_NAME)
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

push-all: build-all
	$(MAKE) OP=push build-all



build-all:
	$(MAKE) ${OP} WHAT=amazon-kernel
	$(MAKE) ${OP} WHAT=amazonlinux   RELEASE=2          IS_LATEST=true
	$(MAKE) ${OP} WHAT=alpine
	$(MAKE) ${OP} WHAT=opensuse      RELEASE=leap       IS_LATEST=true
	$(MAKE) ${OP} WHAT=opensuse      RELEASE=tumbleweed
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=16.04						IS_MANIFEST_LIST=0
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=18.04      				IS_MANIFEST_LIST=1 GOARCH=arm64
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=18.04      				IS_MANIFEST_LIST=1 GOARCH=amd64
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=19.04						IS_MANIFEST_LIST=0
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=20.04      IS_LATEST=true	IS_MANIFEST_LIST=1 GOARCH=arm64
	$(MAKE) ${OP} WHAT=ubuntu        RELEASE=20.04      IS_LATEST=true	IS_MANIFEST_LIST=1 GOARCH=amd64
	$(MAKE) ${OP} WHAT=centos        RELEASE=7
	$(MAKE) ${OP} WHAT=centos        RELEASE=8          IS_LATEST=true
	$(MAKE) ${OP} WHAT=kubeadm
	$(MAKE) build-image-firekube
