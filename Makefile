UID_GID?=$(shell id -u):$(shell id -g)
FIRECRACKER_VERSION:=$(shell cat hack/FIRECRACKER_VERSION)
GO_VERSION=1.12.6
DOCKER_USER=weaveworks
GIT_VERSION:=$(shell hack/ldflags.sh --version-only)
IMAGE_DEV_TAG=dev
IMAGE_TAG:=$(shell hack/ldflags.sh --image-tag-only)
# IS_DIRTY is 1 if the tree state is dirty, otherwise 0
IS_DIRTY:=$(shell echo ${GIT_VERSION} | grep -o dirty | wc -l)
WHAT?=ignite
PROJECT = github.com/weaveworks/ignite
APIS_DIR = ${PROJECT}/pkg/apis
API_DIRS = ${APIS_DIR}/ignite,${APIS_DIR}/ignite/v1alpha1,${APIS_DIR}/ignite/v1alpha2,${APIS_DIR}/meta/v1alpha1
CACHE_DIR = $(shell pwd)/bin/cache
API_DOCS = docs/api/ignite.md docs/api/meta.md
DOCS_PORT = 8000

## Multi-platform-related stuff
GOARCH ?= amd64
QEMUVERSION=v2.9.1

ifeq ($(GOARCH),amd64)
QEMUARCH=amd64
BASEIMAGE=alpine:3.9
ARCH_SUFFIX=
endif
ifeq ($(GOARCH),arm64)
QEMUARCH=aarch64
BASEIMAGE=arm64v8/alpine:3.9
ARCH_SUFFIX=-aarch64
endif

all: binary
binary:
	$(MAKE) shell COMMAND="make bin/$(GOARCH)/${WHAT}"

install: binary
	sudo cp bin/ignite /usr/local/bin

# Make make execute this target although the file already exists.
.PHONY: bin/$(GOARCH)/ignite bin/$(GOARCH)/ignite-spawn bin/$(GOARCH)/ignited
ignite: bin/$(GOARCH)/ignite
ignited: bin/$(GOARCH)/ignited
# Always update the image when ignite-spawn is updated
ignite-spawn: bin/$(GOARCH)/ignite-spawn image
bin/$(GOARCH)/ignite bin/$(GOARCH)/ignited bin/$(GOARCH)/ignite-spawn: bin/%:
	CGO_ENABLED=0 GOARCH=$(GOARCH) go build -mod=vendor -ldflags "$(shell ./hack/ldflags.sh)" -o bin/$* ./cmd/$*
	ln -sf ./$(GOARCH)/$* bin/$*

.PHONY: bin/$(GOARCH)/Dockerfile
image: bin/$(GOARCH)/Dockerfile
bin/$(GOARCH)/Dockerfile: qemu
	sed -e "s|QEMUARCH|$(QEMUARCH)|g;s|BASEIMAGE|$(BASEIMAGE)|g;" Dockerfile > bin/$(GOARCH)/Dockerfile
ifeq ($(GOARCH),amd64)
	# When building for amd64, remove the qemu stuff, it has no part in the amd64 image
	sed -i "/qemu/d" bin/$(GOARCH)/Dockerfile
else
	# Register /usr/bin/qemu-ARCH-static as the handler for non-x86 binaries in the kernel
	docker run --rm --privileged multiarch/qemu-user-static:register --reset
endif
	docker build -t ${DOCKER_USER}/ignite:${IMAGE_DEV_TAG}-$(GOARCH) \
		--build-arg FIRECRACKER_VERSION=${FIRECRACKER_VERSION} \
		--build-arg ARCH_SUFFIX=${ARCH_SUFFIX} bin/$(GOARCH)
ifeq ($(IS_DIRTY),0)
	docker tag ${DOCKER_USER}/ignite:${IMAGE_DEV_TAG}-$(GOARCH) ${DOCKER_USER}/ignite:${IMAGE_TAG}-$(GOARCH)
endif

image-push: image
ifeq ($(IS_DIRTY),0)
	docker push ${DOCKER_USER}/ignite:${IMAGE_TAG}
endif

tidy: $(API_DOCS)
	go mod tidy
	go mod vendor
	hack/generate-client.sh
	gofmt -s -w pkg cmd
	goimports -w pkg cmd
	go run hack/cobra.go

graph:
	hack/graph.sh

.PHONY: $(API_DOCS)
$(API_DOCS): docs/api/%.md: $(CACHE_DIR)/go/bin/godoc2md
	mkdir -p $$(dirname $@) bin/tmp/$*
	mv $(shell pwd)/pkg/apis/$*/v1alpha1/zz_generated* bin/tmp/$*
	$(MAKE) shell COMMAND="/go/bin/godoc2md /go/src/${PROJECT}/pkg/apis/$*/v1alpha1 > $@"
	sed -e "s|src/target|pkg/apis/$*/v1alpha1|g" -i $@
	sed -e "s|(/pkg/apis|(https://github.com/weaveworks/ignite/tree/master/pkg/apis|g" -i $@
	mv bin/tmp/$*/*.go $(shell pwd)/pkg/apis/$*/v1alpha1/
	rm -r bin/tmp/$*

$(CACHE_DIR)/go/bin/godoc2md:
	mkdir -p $(CACHE_DIR)/go/bin/
	curl -sSL https://github.com/luxas/godoc2md/releases/download/v0.1.0/godoc2md > $@
	chmod +x $@

shell:
	mkdir -p $(CACHE_DIR)/go $(CACHE_DIR)/cache
	docker run -it --rm \
		-v $(CACHE_DIR)/go:/go \
		-v $(CACHE_DIR)/cache:/.cache/go-build \
		-v $(shell pwd):/go/src/${PROJECT} \
		-w /go/src/${PROJECT} \
		-u $(shell id -u):$(shell id -g) \
		-e GO111MODULE=on \
		golang:$(GO_VERSION) \
		$(COMMAND)

autogen:
	$(MAKE) shell COMMAND="make dockerized-autogen"

dockerized-autogen: /go/bin/deepcopy-gen /go/bin/defaulter-gen /go/bin/conversion-gen /go/bin/openapi-gen
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
	
	/go/bin/openapi-gen \
		--input-dirs ${API_DIRS} \
		--output-package ${PROJECT}/api/openapi \
		--report-filename api/openapi/violations.txt \
		-h /tmp/boilerplate

/go/bin/%: vendor
	go install k8s.io/code-generator/cmd/$*

/go/bin/openapi-gen:
	go install k8s.io/kube-openapi/cmd/openapi-gen

# QEMU stuff
qemu: bin/$(GOARCH)/qemu-$(QEMUARCH)-static
bin/$(GOARCH)/qemu-$(QEMUARCH)-static:
	mkdir -p bin/$(GOARCH)
ifeq ($(GOARCH),amd64)
	#touch $@
else
	curl -sSL https://github.com/multiarch/qemu-user-static/releases/download/$(QEMUVERSION)/x86_64_qemu-$(QEMUARCH)-static.tar.gz | tar -xz -C bin/$(GOARCH)
	chmod 0755 $@
endif

# Read the docs stuff
build-docs:
	@cd docs && docker build -t ignite-docs .

test-docs: build-docs
	@docker run -it ignite-docs /usr/bin/linkchecker _build/html/index.html

serve-docs: build-docs
	@echo Stating docs website on http://localhost:${DOCS_PORT}/_build/html/index.html
	@docker run -i -p ${DOCS_PORT}:8000 -e USER_ID=$$UID ignite-docs
