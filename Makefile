UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.12
DOCKER_USER=weaveworks
PROJECT = github.com/weaveworks/ignite
APIS_DIR = ${PROJECT}/pkg/apis
CACHE_DIR = /tmp/go-cache

all: binary
binary:
	$(MAKE) shell COMMAND="make bin/ignite"

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

shell:
	mkdir -p $(CACHE_DIR)/bin $(CACHE_DIR)/src $(CACHE_DIR)/cache bin/cache
	docker run -it \
		-v $(CACHE_DIR)/bin:/go/bin \
		-v $(CACHE_DIR)/src:/go/src \
		-v $(CACHE_DIR)/cache:/.cache/go-build \
		-v $(shell pwd):/go/src/github.com/weaveworks/ignite \
		-w /go/src/github.com/weaveworks/ignite \
		-u $(shell id -u):$(shell id -g) \
		-e GO111MODULE=on \
		golang:$(GO_VERSION) \
		$(COMMAND)

autogen:
	$(MAKE) shell COMMAND="make dockerized-autogen"

dockerized-autogen: /go/bin/deepcopy-gen /go/bin/defaulter-gen /go/bin/conversion-gen
	# Let the boilerplate be empty
	touch /tmp/boilerplate
	/go/bin/deepcopy-gen \
		--input-dirs ${APIS_DIR}/ignite/v1alpha1 \
		--bounding-dirs ${APIS_DIR} \
		-O zz_generated.deepcopy \
		-h /tmp/boilerplate 

	/go/bin/defaulter-gen \
		--input-dirs ${APIS_DIR}/ignite/v1alpha1 \
		-O zz_generated.defaults \
		-h /tmp/boilerplate

	/go/bin/conversion-gen \
		--input-dirs ${APIS_DIR}/ignite/v1alpha1 \
		-O zz_generated.conversion \
		-h /tmp/boilerplate

/go/bin/%: vendor
	go install k8s.io/code-generator/cmd/$*
