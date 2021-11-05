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

REGISTRY?=weaveworks
IMAGE_NAME?=${REGISTRY}/ignite-kernel
KERNEL_BUILDER_NAME=${REGISTRY}/ignite-kernel-builder
# Check https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/refs/ for updates
KERNEL_VERSIONS ?= 4.14.254 4.19.215 5.4.157 5.10.77 5.14.16 # If you update this, please keep the .github/workflows/release-kernel-images.yml matrix up to date
GOARCH?=amd64
GOARCH_LIST = amd64 arm64

ifeq ($(GOARCH),amd64)
KERNEL_ARCH=x86
VMLINUX_PATH=vmlinux
endif
ifeq ($(GOARCH),arm64)
ARCH_MAKE_PARAMS="ARCH=arm64 CROSS_COMPILE=aarch64-linux-gnu-"
KERNEL_ARCH=arm64
VMLINUX_PATH=arch/arm64/boot/Image
endif

all: build

kernel-builder:
	$(DOCKER) build -t ${KERNEL_BUILDER_NAME}:dev \
		-f builder/Dockerfile .

fc-upstream:
	curl -L https://raw.githubusercontent.com/firecracker-microvm/firecracker/main/resources/microvm-kernel-x86_64.config > upstream/config-amd64
	curl -L https://raw.githubusercontent.com/firecracker-microvm/firecracker/main/resources/microvm-kernel-arm64.config > upstream/config-arm64

upgrade: $(addprefix upgrade-,$(KERNEL_VERSIONS))

upgrade-%: kernel-builder
	mkdir -p versioned generated
	for upstream in upstream/*; do \
		config="$$(basename $$upstream)-$*"; \
		./upgrade-config.sh $$upstream versioned/$${config}; \
		./patch-config.sh versioned/$${config} generated/$${config} ./config-patches; \
		./upgrade-config.sh generated/$${config} generated/$${config}; \
	done

build: $(addprefix build-,$(KERNEL_VERSIONS))
build-%: kernel-builder
	$(DOCKER) build -t $(IMAGE_NAME):$*-${GOARCH} \
		--build-arg KERNEL_VERSION=$* \
		--build-arg ARCH=${KERNEL_ARCH} \
		--build-arg GOARCH=${GOARCH} \
		--build-arg ARCH_MAKE_PARAMS=${ARCH_MAKE_PARAMS} \
		--build-arg VMLINUX_PATH=${VMLINUX_PATH} .

ctr-import: $(addprefix ctr-import-,$(KERNEL_VERSIONS))
ctr-import-%:
	$(DOCKER) image save $(IMAGE_NAME):$*-${GOARCH} \
		| $(CTR) -n firecracker image import -

push: $(addprefix push-,$(KERNEL_VERSIONS))
push-%:
	../../hack/push-manifest-list.sh $(IMAGE_NAME):$* $(GOARCH_LIST)

pre-release: $(addprefix pre-release-,$(GOARCH_LIST))
pre-release-%:
	$(MAKE) GOARCH=$* build

release: pre-release
	$(MAKE) push
