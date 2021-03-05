FROM BASEIMAGE AS build

# If we're building for another architecture than amd64, this let's us emulate an other platform's docker build.
# If we're building normally, for amd64, this line is removed
COPY qemu-QEMUARCH-static /usr/bin/

# device-mapper is needed for snapshot functionalities
RUN apk add --no-cache \
    device-mapper

# Download the Firecracker binary from Github
ARG FIRECRACKER_VERSION
# If amd64 is set, this is "-x86_64". If arm64, this should be "-aarch64".
ARG FIRECRACKER_ARCH_SUFFIX
RUN wget -qO- https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION}${FIRECRACKER_ARCH_SUFFIX}.tgz | tar -xvz && \
    mv release-${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION}${FIRECRACKER_ARCH_SUFFIX} /usr/local/bin/firecracker && \
    rm -r release-${FIRECRACKER_VERSION}

# Add ignite-spawn to the image
ADD ./ignite-spawn /usr/local/bin/ignite-spawn

# Symlink both firecracker and ignite-spawn to /, too
RUN chmod +x /usr/local/bin/firecracker /usr/local/bin/ignite-spawn && \
    ln -s /usr/local/bin/firecracker  /firecracker  && \
    ln -s /usr/local/bin/ignite-spawn /ignite-spawn

# Create a directory to host any volumes exposed from host
RUN mkdir /volumes

# Use a multi-stage build to allow the resulting image to only consist of one layer
# This makes it more lightweight
FROM scratch
COPY --from=build / /
# ignite-spawn runs as PID 1 in the container, spawning the firecracker process
ENTRYPOINT ["/usr/local/bin/ignite-spawn"]
