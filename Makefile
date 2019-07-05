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
API_DIRS = ${APIS_DIR}/ignite/v1alpha1,${APIS_DIR}/meta/v1alpha1
CACHE_DIR = $(shell pwd)/bin/cache
API_DOCS = api/ignite.md api/meta.md

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
	docker tag ${DOCKER_USER}/ignite:${IMAGE_DEV_TAG} ${DOCKER_USER}/ignite:${IMAGE_TAG}
endif

image-push: image
ifeq ($(IS_DIRTY),0)
	docker push ${DOCKER_USER}/ignite:${IMAGE_TAG}
endif

tidy: $(API_DOCS)
	go mod tidy
	go mod vendor
	gofmt -s -w pkg cmd
	goimports -w pkg cmd
	hack/generate-client.sh
	go run hack/cobra.go

.PHONY: $(API_DOCS)
$(API_DOCS): api/%.md: $(CACHE_DIR)/go/bin/godoc2md
	mkdir -p $$(dirname $@) bin/tmp/$*
	mv $(shell pwd)/pkg/apis/$*/v1alpha1/zz_generated* bin/tmp/$*
	$(MAKE) shell COMMAND="/go/bin/godoc2md /go/src/${PROJECT}/pkg/apis/$*/v1alpha1 > $@"
	mv bin/tmp/$*/*.go $(shell pwd)/pkg/apis/$*/v1alpha1/
	rm -r bin/tmp/$*

$(CACHE_DIR)/go/bin/godoc2md:
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
