#!/bin/bash

KERNEL_BUILDER_IMAGE=luxas/kernel-builder:gcc-7-2

if [[ $# != 2 ]]; then
    echo "Usage: $0 [FROM] [TO]"
    exit 1
fi

FROM=$1
TO=$2
VERSION=$(echo ${TO} | rev | cut -d- -f1 | rev)

if [[ ${FROM} != ${TO} ]]; then
    cp ${FROM} ${TO}
fi

docker run -it \
	-v $(pwd)/${TO}:/tmp/.config \
    ${KERNEL_BUILDER_IMAGE} /bin/bash -c "\
        git checkout v${VERSION} && \
        make clean && make mrproper && cp /tmp/.config . && \
        make EXTRAVERSION="" LOCALVERSION= olddefconfig && \
        cp .config /tmp/.config"
