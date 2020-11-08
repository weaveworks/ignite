FROM weaveworks/ignite-kernel-builder:dev AS builder

# Set up environment variables
ENV CCACHE_DIR=/ccache       \
    SRC_DIR=/usr/src         \
    DIST_DIR=/dist           \
    LINUX_DIR=/usr/src/linux \
    LINUX_REPO_URL=git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

ARG KERNEL_VERSION
ARG KERNEL_EXTRA
# ARCH here is KERNEL_ARCH in the Makefile. It needs to be hardcoded to ARCH for Kconfig to understand
ARG ARCH
ARG GOARCH
ARG ARCH_MAKE_PARAMS

# Clone the desired kernel version and make sure the environment is clean
RUN mkdir -p ${SRC_DIR} ${CCACHE_DIR} ${DIST_DIR} && \
    git clone --depth 1 --branch v${KERNEL_VERSION} ${LINUX_REPO_URL} ${LINUX_DIR} && \
    cd ${LINUX_DIR} && make clean && make mrproper

# Change workdir to run all the build commands from LINUX_DIR.
WORKDIR ${LINUX_DIR}

COPY generated/config-${GOARCH}-${KERNEL_VERSION}${KERNEL_EXTRA} .config

RUN make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= olddefconfig

RUN	make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= -j32
RUN make ${ARCH_MAKE_PARAMS} EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= modules_install

# VMLINUX_PATH is configurable, as the arm64 kernel to be booted by Firecracker is present
# in arch/arm64/boot/Image.
ARG VMLINUX_PATH="vmlinux"
RUN cp ${VMLINUX_PATH} /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} && \
	ln -s /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} /boot/vmlinux && \
	cp .config /boot/config-${KERNEL_VERSION}${KERNEL_EXTRA}

FROM scratch
COPY --from=builder /boot /boot
COPY --from=builder /lib/modules /lib/modules
