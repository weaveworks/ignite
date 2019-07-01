UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.12.6
DOCKER_USER=weaveworks
IMAGE_DEV_TAG=dev
GIT_VERSION:=$(shell hack/ldflags.sh --version-only)
# IS_DIRTY is 1 if the tree state is dirty, otherwise 0
IS_DIRTY:=$(shell echo ${GIT_VERSION} | grep -o dirty | wc -l)
WHAT?=ignite


all: binary
binary:
	docker run -it --rm -v $(shell pwd):/build -w /build golang:${GO_VERSION} sh -c "\
		make ${WHAT} && \
		chown ${UID_GID} bin/${WHAT}"

install: binary
	sudo cp bin/ignite /usr/local/bin

# Make make execute this target although the file already exists.
.PHONY: bin/ignite bin/ignite-spawn
ignite: bin/ignite
# Always update the image when ignite-spawn is updated
ignite-spawn: bin/ignite-spawn image
bin/ignite bin/ignite-spawn: bin/%:
	CGO_ENABLED=0 go build -mod=vendor -ldflags "$(shell ./hack/ldflags.sh)" -o bin/$* ./cmd/$*

image:
	docker build -t ${DOCKER_USER}/ignite:${IMAGE_DEV_TAG} \
		--build-arg FIRECRACKER_VERSION=${FIRECRACKER_VERSION} .
ifeq ($(IS_DIRTY),0)
	docker tag ${DOCKER_USER}/ignite:${IMAGE_DEV_TAG} ${DOCKER_USER}/ignite:${GIT_VERSION}
endif

image-push: image
ifeq ($(IS_DIRTY),0)
	docker push ${DOCKER_USER}/ignite:${GIT_VERSION}
endif

tidy:
	go mod tidy
	go mod vendor
	gofmt -s -w pkg cmd
	goimports -w pkg cmd
	go run hack/cobra.go
