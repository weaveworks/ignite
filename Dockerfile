FROM alpine:3.9

RUN apk add --no-cache \
    iproute2 \
    device-mapper

ARG FIRECRACKER_VERSION
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /ignite/firecracker

# The ignite binary should be bind-mounted over /ignite/ignite
# The downloaded Firecracker binary exists in /ignite/firecracker,
# but both fc and ignite are symlinked to be in $PATH
# The data directory is mounted in from the host to /var/lib/firecracker
RUN touch /ignite/ignite && \
    chmod +x /ignite/firecracker && \
    ln -s /ignite/firecracker /usr/local/bin/firecracker && \
    ln -s /ignite/ignite /usr/local/bin/ignite

# Ignite runs as PID in the container, spawning the firecracker process
ENTRYPOINT ["/ignite/ignite", "container"]
