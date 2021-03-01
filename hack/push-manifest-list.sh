#!/usr/bin/env bash

if [[ $# == 0 ]]; then
    echo "Usage: $0 [manifest-list] [arch0...n]"
    exit 1
fi

set -x

IMAGE=${1}
shift
ARCHES="$@"
DOCKER=${DOCKER:-docker}

export DOCKER_CLI_EXPERIMENTAL=enabled

for arch in ${ARCHES}; do
    ${DOCKER} push ${IMAGE}-${arch}
done

sleep 1

${DOCKER} manifest create --amend ${IMAGE} $(echo ${ARCHES} | sed -e "s~[^ ]*~${IMAGE}\-&~g")
for arch in ${ARCHES}; do
    ${DOCKER} manifest annotate --arch=${arch} ${IMAGE} ${IMAGE}-${arch}
done
${DOCKER} manifest push --purge ${IMAGE}

for arch in ${ARCHES}; do
    CURR_IMAGE=${IMAGE}-${arch}
    ${DOCKER} manifest create --amend ${CURR_IMAGE} ${CURR_IMAGE}
    ${DOCKER} manifest annotate --arch=${arch} ${CURR_IMAGE} ${CURR_IMAGE}
    ${DOCKER} manifest push --purge ${CURR_IMAGE}
done

