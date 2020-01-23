FROM luxas/kernel-builder:gcc-7 AS builder

# Install crosscompilers for non-amd64 arches
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        binutils-multiarch \
        gcc-7-aarch64-linux-gnu && \
	ln -s /usr/bin/aarch64-linux-gnu-gcc-7 /usr/bin/aarch64-linux-gnu-gcc

# Fetch the latest tags from the linux kernel repo 
RUN git fetch --tags

ARG KERNEL_VERSION
ARG KERNEL_EXTRA
# ARCH here is KERNEL_ARCH in the Makefile. It needs to be hardcoded to ARCH for Kconfig to understand
ARG ARCH
ARG GOARCH
ARG ARCH_MAKE_PARAMS

# Checkout the desired kernel version and make sure the environment is clean
RUN git checkout v${KERNEL_VERSION} && \
    make clean && make mrproper

COPY generated/config-${GOARCH}-${KERNEL_VERSION}${KERNEL_EXTRA} .config

RUN make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= olddefconfig

RUN	make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= -j32
RUN make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= modules_install

RUN cp vmlinux /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} && \
	ln -s /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} /boot/vmlinux && \
	cp .config /boot/config-${KERNEL_VERSION}${KERNEL_EXTRA}

FROM scratch
COPY --from=builder /boot /boot
COPY --from=builder /lib/modules /lib/modules
