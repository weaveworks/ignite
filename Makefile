UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.12.6
DOCKER_USER=weaveworks
IMAGE_DEV_TAG=dev
GIT_VERSION:=$(shell hack/ldflags.sh --version-only)
# IS_DIRTY is 1 if the tree state is dirty, otherwise 0
IS_DIRTY:=$(shell echo ${GIT_VERSION} | grep -o dirty | wc -l)
WHAT?=ignite
PROJECT = github.com/weaveworks/ignite
APIS_DIR = ${PROJECT}/pkg/apis
API_DIRS = ${APIS_DIR}/ignite/v1alpha1,${APIS_DIR}/meta/v1alpha1
CACHE_DIR = /tmp/go-cache

all: binary
binary:
	$(MAKE) shell COMMAND="make bin/${WHAT}"

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
		--input-dirs ${API_DIRS} \
		--bounding-dirs ${APIS_DIR} \
		-O zz_generated.deepcopy \
		-h /tmp/boilerplate 

	/go/bin/defaulter-gen \
		--input-dirs ${API_DIRS} \
		-O zz_generated.defaults \
		-h /tmp/boilerplate

	/go/bin/conversion-gen \
		--input-dirs ${API_DIRS} \
		-O zz_generated.conversion \
		-h /tmp/boilerplate

/go/bin/%: vendor
	go install k8s.io/code-generator/cmd/$*
