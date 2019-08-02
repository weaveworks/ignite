FROM BASEIMAGE AS build

# If we're building for another architecture than amd64, this let's us emulate an other platform's docker build.
# If we're building normally, for amd64, this line is removed
COPY qemu-QEMUARCH-static /usr/bin/

# Install iproute2 for access to the "ip" command. In the future this dependency will be removed
# device-mapper is needed for the snapshot functionalities
RUN apk add --no-cache \
    iproute2 \
    device-mapper

# Download the Firecracker binary from Github
ARG FIRECRACKER_VERSION
# If amd64 is set, this is "". If arm64, this should be "-aarch64".
ARG ARCH_SUFFIX
RUN wget -q -O /usr/local/bin/firecracker https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION}${ARCH_SUFFIX}

# Add ignite-spawn to the image
ADD ./ignite-spawn /usr/local/bin/ignite-spawn

# Symlink both firecracker and ignite-spawn to /, too
RUN chmod +x /usr/local/bin/firecracker /usr/local/bin/ignite-spawn && \
    ln -s /usr/local/bin/firecracker  /firecracker  && \
    ln -s /usr/local/bin/ignite-spawn /ignite-spawn

# Use a multi-stage build to allow the resulting image to only consist of one layer
# This makes it more lightweight
FROM scratch
COPY --from=build / /
# ignite-spawn runs as PID 1 in the container, spawning the firecracker process
ENTRYPOINT ["/usr/local/bin/ignite-spawn"]
