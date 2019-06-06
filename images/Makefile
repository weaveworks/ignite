# WHAT specifies the OS image to build
WHAT?=
DOCKER_USER?=weaveworks
VERSION?=$(shell git describe HEAD --tags)

all: build
build:
ifeq ($(WHAT),)
$(error WHAT is a required argument)
endif
	@ls ${WHAT} >/dev/null
	docker build -t ${DOCKER_USER}/ignite-${WHAT}:${VERSION} ${WHAT}
