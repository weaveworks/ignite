#!/usr/bin/env bash

set -eu

# This script runs two local, private, docker registries,
# one with with self-signed certificates, TLS, and HTTP,
# the other with plain HTTP, both using basic auth.

HTTP_REGISTRY_SECRET_PATH="$(mktemp -d)/ignite-test-registry-http"
HTTPS_REGISTRY_SECRET_PATH="$(mktemp -d)/ignite-test-registry-https"

PRIVATE_KEY="${HTTPS_REGISTRY_SECRET_PATH}/certs/domain.key"
CERT="${HTTPS_REGISTRY_SECRET_PATH}/certs/domain.crt"

HTTP_USERNAME="http_testuser"
HTTP_PASSWORD="http_testpassword"
HTTPS_USERNAME="https_testuser"
HTTPS_PASSWORD="https_testpassword"

BIND_IP="127.5.0.1"
HTTP_ADDR="${BIND_IP}:5080"
HTTPS_ADDR="${BIND_IP}:5443"

OS_IMG="weaveworks/ignite-ubuntu:latest"
KERNEL_IMG="weaveworks/ignite-kernel:5.4.108"
HTTP_LOCAL_OS_IMG="${HTTP_ADDR}/weaveworks/ignite-ubuntu:test"
HTTP_LOCAL_KERNEL_IMG="${HTTP_ADDR}/weaveworks/ignite-kernel:test"
HTTPS_LOCAL_OS_IMG="${HTTPS_ADDR}/weaveworks/ignite-ubuntu:test"
HTTPS_LOCAL_KERNEL_IMG="${HTTPS_ADDR}/weaveworks/ignite-kernel:test"

# Clear any existing registry secret and create new directories.
mkdir -p "${HTTP_REGISTRY_SECRET_PATH}/auth"
mkdir -p "${HTTPS_REGISTRY_SECRET_PATH}/"{certs,auth}

# Generate key and cert.
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=US/ST=Foo/L=Bar/O=Weave" \
    -keyout "${PRIVATE_KEY}" -out "${CERT}"
chmod 400 "${PRIVATE_KEY}"

# Create htpasswd files.
docker run --rm \
  --entrypoint htpasswd \
  httpd:2 -Bbn "${HTTP_USERNAME}" "${HTTP_PASSWORD}" > "${HTTP_REGISTRY_SECRET_PATH}/auth/htpasswd"

docker run --rm \
  --entrypoint htpasswd \
  httpd:2 -Bbn "${HTTPS_USERNAME}" "${HTTPS_PASSWORD}" > "${HTTPS_REGISTRY_SECRET_PATH}/auth/htpasswd"

# Run the registries
docker run -d --rm \
  --name ignite-test-http-registry \
  -v "${HTTP_REGISTRY_SECRET_PATH}/auth":/auth \
  -e REGISTRY_AUTH=htpasswd \
  -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -p "${HTTP_ADDR}":5000 \
  registry:2

docker run -d --rm \
  --name ignite-test-https-registry \
  -v "${HTTPS_REGISTRY_SECRET_PATH}/auth":/auth \
  -v "${HTTPS_REGISTRY_SECRET_PATH}/certs":/certs \
  -e REGISTRY_AUTH=htpasswd \
  -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
  -p "${HTTPS_ADDR}":5000 \
  registry:2

# Login, push test images to download in tests and logout.
docker pull "${OS_IMG}"
docker pull "${KERNEL_IMG}"

docker tag "${OS_IMG}" "${HTTP_LOCAL_OS_IMG}"
docker tag "${KERNEL_IMG}" "${HTTP_LOCAL_KERNEL_IMG}"
docker tag "${OS_IMG}" "${HTTPS_LOCAL_OS_IMG}"
docker tag "${KERNEL_IMG}" "${HTTPS_LOCAL_KERNEL_IMG}"

docker login -u "${HTTP_USERNAME}" -p "${HTTP_PASSWORD}" "https://${HTTP_ADDR}"
docker login -u "${HTTPS_USERNAME}" -p "${HTTPS_PASSWORD}" "https://${HTTPS_ADDR}"

# push in parallel, block until all finished
docker push "${HTTP_LOCAL_OS_IMG}" &
docker push "${HTTP_LOCAL_KERNEL_IMG}" &
docker push "${HTTPS_LOCAL_OS_IMG}" &
docker push "${HTTPS_LOCAL_KERNEL_IMG}" &
wait

docker logout "http://${HTTP_ADDR}"
docker logout "https://${HTTPS_ADDR}"

docker rmi "${HTTP_LOCAL_OS_IMG}" "${HTTP_LOCAL_KERNEL_IMG}"
docker rmi "${HTTPS_LOCAL_OS_IMG}" "${HTTPS_LOCAL_KERNEL_IMG}"
