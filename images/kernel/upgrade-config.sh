#!/bin/bash

KERNEL_BUILDER_IMAGE=luxas/kernel-builder:gcc-7

if [[ $# != 2 ]]; then
    echo "Usage: $0 [FROM] [TO]"
    exit 1
fi

FROM=$1
TO=$2
VERSION="$(echo ${TO} | rev | cut -d- -f1 | rev)"  # Extracts the trailing hyphenated field -- this is dependent on naming the resulting file properly (ex: ./versioned/config-amd64-5.4.13)
ARCH=$(echo ${TO} | cut -d- -f2)

# Only set the extra flag for non-amd64 arches
if [[ ${ARCH} != amd64 ]]; then
    ARCH_PARAMETER="-e ARCH=${ARCH}"
fi

if [[ ${FROM} != ${TO} ]]; then
    cp ${FROM} ${TO}
fi

if ! [ -x "$(command -v podman)" ]; then
  oci_tool=podman
else
  oci_tool=docker
fi

${oci_tool} run -it \
    ${ARCH_PARAMETER} \
	-v $(pwd)/${TO}:/tmp/.config \
    ${KERNEL_BUILDER_IMAGE} /bin/bash -c "\
        git checkout v${VERSION} && \
        make clean && make mrproper && cp /tmp/.config . && \
        make EXTRAVERSION="" LOCALVERSION= olddefconfig && \
        cp .config /tmp/.config"
