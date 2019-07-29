FROM alpine:3.9 AS build

# Install iproute2 for access to the "ip" command. In the future this dependency will be removed
# device-mapper is needed for the snapshot functionalities
RUN apk add --no-cache \
    iproute2 \
    device-mapper

# Download the Firecracker binary from Github
ARG FIRECRACKER_VERSION
RUN wget -q -O /usr/local/bin/firecracker https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION}

# Add ignite-spawn to the image
ADD bin/ignite-spawn /usr/local/bin/ignite-spawn

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
