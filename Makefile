MIN_MAKE_VERSION = 3.82
ifneq ($(MIN_MAKE_VERSION), $(firstword $(sort $(MAKE_VERSION) $(MIN_MAKE_VERSION))))
$(error this project requires make version $(MIN_MAKE_VERSION) or higher)
endif

SHELL:=/bin/bash
# Set the command for running `docker`
# -- allows user to override for things like sudo usage or container images
DOCKER := docker
# DOCKER_TTY defaults for a nice docker buildkit tty experience on dev laptops -- can be disabled in CI if no tty available
DOCKER_TTY ?= -t

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
UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.17.9
DOCKER_USER?=weaveworks
IMAGE=$(DOCKER_USER)/ignite
GIT_VERSION:=$(shell DOCKER_USER=$(DOCKER_USER) hack/ldflags.sh --version-only)
IMAGE_DEV_TAG=dev
IMAGE_TAG:=$(shell IGNITE_GIT_VERSION=$(GIT_VERSION) DOCKER_USER=$(DOCKER_USER) hack/ldflags.sh --image-tag-only)
# IS_DIRTY is 1 if the tree state is dirty, otherwise 0
IS_DIRTY:=$(shell echo ${GIT_VERSION} | grep -c dirty)
PROJECT = github.com/weaveworks/ignite
APIS_DIR = ${PROJECT}/pkg/apis
API_DIRS = ${APIS_DIR}/ignite,${APIS_DIR}/ignite/v1alpha2,${APIS_DIR}/ignite/v1alpha3,${APIS_DIR}/ignite/v1alpha4,${APIS_DIR}/meta/v1alpha1
CACHE_DIR = $(shell pwd)/bin/cache
# Specifies if this is a CI build or not; if it is, it will save the docker image created to bin/$(GOARCH)/image.tar
IS_CI_BUILD ?= 0
# Unset any GOFLAGS that would interfere with the build
undefine GOFLAGS

## Multi-platform-related stuff
# Ignore GOARCH from the environment for the query, otherwise it results in the environment's GOARCH
GOHOSTARCH := $(shell GOARCH= go env GOARCH 2>/dev/null || echo "amd64")
GOARCH ?= amd64
GOARCH_LIST = amd64 arm64
QEMUVERSION=v4.2.0-6

ifeq ($(GOARCH),amd64)
QEMUARCH=amd64
BASEIMAGE=alpine:3.13
FIRECRACKER_ARCH_SUFFIX=-x86_64
endif
ifeq ($(GOARCH),arm64)
QEMUARCH=aarch64
BASEIMAGE=arm64v8/alpine:3.13
FIRECRACKER_ARCH_SUFFIX=-aarch64
endif

E2E_REGEX := Test
E2E_COUNT := 1

TEST_COUNT := 1

# Project packages.
TEST_REQUIRES_ROOT_PACKAGES="$(PROJECT)/pkg/runtime/containerd"
TEST_E2E_PACKAGES="$(PROJECT)/e2e"


# Default is to build all the binaries for this architecture
all: build-all-$(GOARCH)

install: ignite
	sudo cp bin/$(GOARCH)/ignite /usr/local/bin

install-all: install ignited
	sudo cp bin/$(GOARCH)/ignited /usr/local/bin

BINARIES = ignite ignited ignite-spawn
$(BINARIES):
	$(MAKE) go-make TARGETS="bin/$(GOARCH)/$@"
	# Always update the image when ignite-spawn is updated
	[[ $@ == "ignite-spawn" ]] && $(MAKE) image || exit 0

GO_MAKE_TARGET := go-in-docker # user can override to "local"
go-make:
# MAKEFLAGS must be set manually in case the sub-make is in a container:
	$(MAKE) $(GO_MAKE_TARGET) COMMAND="$(MAKE) \"MAKEFLAGS=$(MAKEFLAGS)\" $(TARGETS)"
local: # Do not use directly -- use $(GO_MAKE_TARGET)
	$(COMMAND)
go-in-docker: # Do not use directly -- use $(GO_MAKE_TARGET)
	mkdir -p $(CACHE_DIR)/go $(CACHE_DIR)/cache
	$(DOCKER) run -i $(DOCKER_TTY) --rm \
		-v $(CACHE_DIR)/go:/go \
		-v $(CACHE_DIR)/cache:/.cache/go-build \
		-v $(shell pwd):/go/src/${PROJECT} \
		-w /go/src/${PROJECT} \
		-u $(shell id -u):$(shell id -g) \
		-e GOARCH=$(GOARCH) \
		golang:$(GO_VERSION) \
		$(COMMAND)

# Make make execute this target although the file already exists.
.PHONY: bin/$(GOARCH)/ignite bin/$(GOARCH)/ignite-spawn bin/$(GOARCH)/ignited
bin/$(GOARCH)/ignite bin/$(GOARCH)/ignited bin/$(GOARCH)/ignite-spawn: bin/$(GOARCH)/%:
	CGO_ENABLED=0 GOARCH=$(GOARCH) go build -mod=vendor -ldflags "$(shell IGNITE_GIT_VERSION=$(GIT_VERSION) DOCKER_USER=$(DOCKER_USER) ./hack/ldflags.sh)" -o bin/$(GOARCH)/$* ./cmd/$*
ifeq ($(GOARCH),$(GOHOSTARCH))
	ln -sf ./$(GOARCH)/$* bin/$*
endif

.PHONY: bin/$(GOARCH)/Dockerfile
image: bin/$(GOARCH)/Dockerfile
bin/$(GOARCH)/Dockerfile: qemu
	sed -e "s|QEMUARCH|$(QEMUARCH)|g;s|BASEIMAGE|$(BASEIMAGE)|g;" Dockerfile > bin/$(GOARCH)/Dockerfile
ifeq ($(GOARCH),amd64)
	# When building for amd64, remove the qemu stuff, it has no part in the amd64 image
	sed -i "/qemu/d" bin/$(GOARCH)/Dockerfile
else
	# Register /usr/bin/qemu-ARCH-static as the handler for non-x86 binaries in the kernel
	$(DOCKER) run --rm --privileged multiarch/qemu-user-static:register --reset
endif
	$(DOCKER) build -t $(IMAGE):${IMAGE_DEV_TAG}-$(GOARCH) \
		--build-arg FIRECRACKER_VERSION=${FIRECRACKER_VERSION} \
		--build-arg FIRECRACKER_ARCH_SUFFIX=${FIRECRACKER_ARCH_SUFFIX} bin/$(GOARCH)
	$(DOCKER) image save $(IMAGE):${IMAGE_DEV_TAG}-$(GOARCH) \
		| $(CTR) -n firecracker image import -
ifeq ($(GOARCH),$(GOHOSTARCH))
	# Only tag the development image if its architecture matches the host
	$(DOCKER) tag $(IMAGE):${IMAGE_DEV_TAG}-$(GOARCH) $(IMAGE):${IMAGE_DEV_TAG}
	$(DOCKER) image save $(IMAGE):${IMAGE_DEV_TAG} \
		| $(CTR) -n firecracker image import -
endif
ifneq (${IMAGE_TAG},)
ifeq ($(IS_DIRTY),0)
	$(DOCKER) tag $(IMAGE):${IMAGE_DEV_TAG}-$(GOARCH) $(IMAGE):${IMAGE_TAG}-$(GOARCH)
	$(DOCKER) image save $(IMAGE):${IMAGE_TAG}-$(GOARCH) \
		| $(CTR) -n firecracker image import -
ifeq ($(GOARCH),$(GOHOSTARCH))
	# For dev builds for a clean (non-dirty) environment; "simulate" that
	# a manifest list has been built by tagging the docker image
	$(DOCKER) tag $(IMAGE):${IMAGE_TAG}-$(GOARCH) $(IMAGE):${IMAGE_TAG}
	$(DOCKER) image save $(IMAGE):${IMAGE_TAG} \
		| $(CTR) -n firecracker image import -
endif
endif
ifeq ($(IS_CI_BUILD),1)
	$(DOCKER) save $(IMAGE):${IMAGE_TAG}-$(GOARCH) -o bin/$(GOARCH)/image.tar
endif
endif

build-all: $(addprefix build-all-,$(GOARCH_LIST))
build-all-%:
	$(MAKE) GOARCH=$* $(BINARIES)

release: build-all
ifneq ($(IS_DIRTY),0)
	$(error "cannot release dirty tree")
endif
	mkdir -p bin/releases/${GIT_VERSION}
	cp -r bin/{amd64,arm64} bin/releases/${GIT_VERSION}
	DOCKER=$(DOCKER) hack/push-manifest-list.sh $(IMAGE):$(IMAGE_TAG) $(GOARCH_LIST)

tidy: /go/bin/goimports
	go mod tidy
	go mod vendor
	hack/generate-client.sh
	gofmt -s -w pkg cmd
	goimports -w pkg cmd
	go run hack/cobra.go

tidy-in-docker:
	$(MAKE) go-make TARGETS="tidy"

# Exclude e2e tests, incomplete package pkg/dm and packages that require root.
test:
	go test -mod=vendor -v $(shell go list ./... | grep -v /e2e | grep -v /pkg/dm | grep -v $(TEST_REQUIRES_ROOT_PACKAGES)) -count=$(TEST_COUNT)

test-in-docker:
	$(MAKE) go-make TARGETS="test"

# Run tests that require root.
root-test:
	sudo -E $(shell which go) test -mod=vendor -v $(TEST_REQUIRES_ROOT_PACKAGES) -count=$(TEST_COUNT)

graph:
	hack/graph.sh

api-docs: godoc2md
	cd pkg/apis/ && for gv in */v1*; do \
		GROUPVERSION=$$gv GROUP_VERSION=$$(echo $$gv | sed 's|/|_|g') \
		$(MAKE) -C ../../ api-doc; done

api-doc:
	mkdir -p docs/api bin/tmp/${GROUPVERSION}
	mv $(shell pwd)/pkg/apis/${GROUPVERSION}/zz_generated* bin/tmp/${GROUPVERSION}
	$(MAKE) $(GO_MAKE_TARGET) COMMAND="godoc2md /go/src/${PROJECT}/pkg/apis/${GROUPVERSION} > bin/tmp/${GROUP_VERSION}.md"
	sed -e "s|src/target|pkg/apis/${GROUPVERSION}|g;s|/go/src/||g" -i bin/tmp/${GROUP_VERSION}.md
	sed -e "s|(/pkg/apis|(https://github.com/weaveworks/ignite/tree/main/pkg/apis|g" -i bin/tmp/${GROUP_VERSION}.md
	mv bin/tmp/${GROUPVERSION}/*.go $(shell pwd)/pkg/apis/${GROUPVERSION}/
	rm -r bin/tmp/${GROUPVERSION}
	# Format the docs with pandoc
	$(DOCKER) run -i $(DOCKER_TTY) --rm \
		-v $(shell pwd):/data \
		-u $(shell id -u):$(shell id -g) \
		pandoc/core:2.10 `# TODO: let's try bumping this image tag in 2021` \
			--from markdown \
			--to gfm \
			bin/tmp/${GROUP_VERSION}.md > docs/api/${GROUP_VERSION}.md

autogen:
	$(MAKE) go-make TARGETS="go-autogen"

go-autogen:
	# Let the boilerplate be empty
	touch /tmp/boilerplate

	go run -mod=vendor ./vendor/k8s.io/code-generator/cmd/deepcopy-gen \
		--input-dirs ${API_DIRS} \
		--bounding-dirs ${APIS_DIR} \
		-O zz_generated.deepcopy \
		-h /tmp/boilerplate 

	go run -mod=vendor ./vendor/k8s.io/code-generator/cmd/defaulter-gen \
		--input-dirs ${API_DIRS} \
		-O zz_generated.defaults \
		-h /tmp/boilerplate

	go run -mod=vendor ./vendor/k8s.io/code-generator/cmd/conversion-gen \
		--input-dirs ${API_DIRS} \
		-O zz_generated.conversion \
		-h /tmp/boilerplate

	go run -mod=vendor ./vendor/k8s.io/kube-openapi/cmd/openapi-gen \
		--input-dirs ${API_DIRS} \
		--output-package ${PROJECT}/pkg/openapi \
		--report-filename pkg/openapi/violations.txt \
		-h /tmp/boilerplate
	
	# These commands modify the env a bit, clean up after us too
	$(MAKE) tidy

godoc2md: bin/cache/go/bin/godoc2md
bin/cache/go/bin/godoc2md:
	mkdir -p $$(dirname $@)
	curl -sSL https://github.com/luxas/godoc2md/releases/download/v0.1.0/godoc2md > $@
	chmod +x $@

/go/bin/goimports:
	go get golang.org/x/tools/cmd/goimports

# QEMU stuff
qemu: bin/$(GOARCH)/qemu-$(QEMUARCH)-static
bin/$(GOARCH)/qemu-$(QEMUARCH)-static:
	mkdir -p bin/$(GOARCH)
ifneq ($(GOARCH),amd64)
	curl -sSL https://github.com/multiarch/qemu-user-static/releases/download/$(QEMUVERSION)/x86_64_qemu-$(QEMUARCH)-static.tar.gz | tar -xz -C bin/$(GOARCH)
	chmod 0755 $@
endif

# Read the docs stuff
bin/docs/builder-image.tar:
	mkdir -p bin/docs
	$(DOCKER) build -t ignite-docs-builder -f docs/Dockerfile.build docs
	$(DOCKER) save ignite-docs-builder -o $@

e2e: build-all e2e-nobuild

e2e-nobuild:
	sudo -E IGNITE_E2E_HOME=$(shell pwd) \
		$(shell which go) test \
		$(TEST_E2E_PACKAGES) -v -mod=vendor \
		-count $(E2E_COUNT) \
		-run $(E2E_REGEX)

docs-deps:
	pip3 install -r docs/requirements.txt

serve-docs: docs-deps
	mkdocs serve
