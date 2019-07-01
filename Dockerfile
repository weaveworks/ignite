FROM alpine:3.9

RUN apk add --no-cache \
    iproute2 \
    device-mapper

ARG FIRECRACKER_VERSION
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /usr/local/bin/firecracker

ADD bin/ignite-spawn /usr/local/bin/ignite-spawn

# Symlink both firecracker and ignite-spawn to /, too
RUN chmod +x /usr/local/bin/firecracker /usr/local/bin/ignite-spawn && \
    ln -s /usr/local/bin/firecracker  /firecracker  && \
    ln -s /usr/local/bin/ignite-spawn /ignite-spawn

# ignite-spawn runs as PID 1 in the container, spawning the firecracker process
ENTRYPOINT ["/usr/local/bin/ignite-spawn"]
