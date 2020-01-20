# Check https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/refs/ for updates
KERNEL_VERSIONS ?= 4.14.166 4.19.97 5.4.13

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
	docker build -t weaveworks/ignite-kernel:$* \
		--build-arg KERNEL_VERSION=$* .

push: build $(addprefix push-,$(KERNEL_VERSIONS))
push-%:
	docker push weaveworks/ignite-kernel:${KERNEL_VERSION}
