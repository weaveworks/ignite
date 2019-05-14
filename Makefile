all: build
build:
	docker build -t ignite .

binary:
	docker run -it -v $(shell pwd):/build -w /build golang:1.12 make ignite

ignite: bin/ignite
bin/ignite:
	CGO_ENABLED=0 go build -mod=vendor -o bin/ignite ./cmd/ignite
