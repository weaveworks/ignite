UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.12
DOCKER_USER=weaveworks

all: binary
binary:
	docker run -it --rm -v $(shell pwd):/build -w /build golang:${GO_VERSION} sh -c "\
		make ignite && \
		chown ${UID_GID} bin/ignite"

install: binary
	sudo cp bin/ignite /usr/local/bin

# Make make execute this target although the file already exists.
.PHONY: bin/ignite
ignite: bin/ignite
bin/ignite:
	CGO_ENABLED=0 go build -mod=vendor -ldflags "$(shell ./hack/ldflags.sh)" -o bin/ignite ./cmd/ignite

image:
	docker build -t ${DOCKER_USER}/ignite:${FIRECRACKER_VERSION} \
		--build-arg FIRECRACKER_VERSION=${FIRECRACKER_VERSION} .

image-push: image
	docker push ${DOCKER_USER}/ignite:${FIRECRACKER_VERSION}

tidy:
	go mod tidy
	go mod vendor
	gofmt -s -w pkg cmd
	goimports -w pkg cmd
	go run hack/cobra.go
