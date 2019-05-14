FROM golang:1.11 AS build
WORKDIR /build
ADD . .
RUN make ignite

FROM alpine:latest
RUN apk add --update iproute2
ARG FIRECRACKER_VERSION=v0.15.2
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /firecracker
# This Dockerfile's context is root of this repo
COPY --from=build /build/bin/ignite /
RUN chmod +x /firecracker /ignite
RUN ln -s /firecracker /usr/local/bin/firecracker && ln -s /ignite /usr/local/bin/ignite
