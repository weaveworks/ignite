FROM alpine:latest

RUN apk add --no-cache iproute2

VOLUME /var/lib/firecracker

ARG FIRECRACKER_VERSION=v0.16.0
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker /ignite/firecracker

# The ignite binary should be bind-mounted over /ignite/ignite
RUN touch /ignite/ignite && \
    chmod +x /ignite/firecracker && \
    ln -s /ignite/firecracker /usr/local/bin/firecracker && \
    ln -s /ignite/ignite /usr/local/bin/ignite

ENTRYPOINT ["/ignite/ignite", "container"]