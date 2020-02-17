FROM luxas/kernel-builder:gcc-7 AS builder
ARG KERNEL_VERSION
ARG KERNEL_EXTRA
RUN git fetch --tags
RUN git checkout v${KERNEL_VERSION} && \
    make clean && make mrproper

COPY generated/config-amd64-${KERNEL_VERSION}${KERNEL_EXTRA} .config

RUN make EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= olddefconfig

RUN	make EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= -j32
RUN make EXTRAVERSION=${KERNEL_EXTRA} LOCALVERSION= modules_install

RUN cp vmlinux /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} && \
	ln -s /boot/vmlinux-${KERNEL_VERSION}${KERNEL_EXTRA} /boot/vmlinux && \
	cp .config /boot/config-${KERNEL_VERSION}${KERNEL_EXTRA}

FROM scratch
COPY --from=builder /boot /boot
COPY --from=builder /lib/modules /lib/modules
