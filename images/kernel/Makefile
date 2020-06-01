# Check https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/refs/ for updates
REGISTRY?=weaveworks
IMAGE_NAME?=${REGISTRY}/ignite-kernel
KERNEL_VERSIONS ?= 4.14.182 4.19.125 5.4.43
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

upgrade: $(addprefix upgrade-,$(KERNEL_VERSIONS))
	./patch-config.sh
	for file in generated/*; do \
		./upgrade-config.sh $$file $$file; done

upgrade-%:
	for file in upstream/*; do \
		./upgrade-config.sh $$file versioned/$$(basename $$file)-$*; done

build: $(addprefix build-,$(KERNEL_VERSIONS))
build-%:
	docker build -t $(IMAGE_NAME):$*-${GOARCH} \
		--build-arg KERNEL_VERSION=$* \
		--build-arg ARCH=${KERNEL_ARCH} \
		--build-arg GOARCH=${GOARCH} \
		--build-arg ARCH_MAKE_PARAMS=${ARCH_MAKE_PARAMS} \
		--build-arg VMLINUX_PATH=${VMLINUX_PATH} .

push: $(addprefix push-,$(KERNEL_VERSIONS))
push-%:
	../../hack/push-manifest-list.sh $(IMAGE_NAME):$* $(GOARCH_LIST)

pre-release: $(addprefix pre-release-,$(GOARCH_LIST))
pre-release-%:
	$(MAKE) GOARCH=$* build

release: pre-release
	$(MAKE) push
