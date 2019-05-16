UID_GID?=$(shell id -u):$(shell id -g)

all: build
build:
	docker build -t ignite .

binary:
	docker run -it --rm -v $(shell pwd):/build -w /build golang:1.12 sh -c "\
		make ignite && \
		chown ${UID_GID} bin/ignite"

ignite: bin/ignite
bin/ignite:
	CGO_ENABLED=0 go build -mod=vendor -o bin/ignite ./cmd/ignite

.PHONY: bin/ignite