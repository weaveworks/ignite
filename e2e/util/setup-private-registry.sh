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
KERNEL_IMG="weaveworks/ignite-kernel:5.10.51"
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

# Use a test config directory to avoid modifying the user's default docker
# configuration.
DOCKER_CONFIG_DIR="$(mktemp -d)/ignite-docker-config"
mkdir -p "${DOCKER_CONFIG_DIR}"

DOCKER="docker --config=${DOCKER_CONFIG_DIR}"

# Create htpasswd files.
${DOCKER} run --rm \
  --entrypoint htpasswd \
  httpd:2 -Bbn "${HTTP_USERNAME}" "${HTTP_PASSWORD}" > "${HTTP_REGISTRY_SECRET_PATH}/auth/htpasswd"

${DOCKER} run --rm \
  --entrypoint htpasswd \
  httpd:2 -Bbn "${HTTPS_USERNAME}" "${HTTPS_PASSWORD}" > "${HTTPS_REGISTRY_SECRET_PATH}/auth/htpasswd"

# Run the registries
${DOCKER} run -d --rm \
  --name ignite-test-http-registry \
  -v "${HTTP_REGISTRY_SECRET_PATH}/auth":/auth \
  -e REGISTRY_AUTH=htpasswd \
  -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -p "${HTTP_ADDR}":5000 \
  registry:2

${DOCKER} run -d --rm \
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
${DOCKER} pull "${OS_IMG}"
${DOCKER} pull "${KERNEL_IMG}"

${DOCKER} tag "${OS_IMG}" "${HTTP_LOCAL_OS_IMG}"
${DOCKER} tag "${KERNEL_IMG}" "${HTTP_LOCAL_KERNEL_IMG}"
${DOCKER} tag "${OS_IMG}" "${HTTPS_LOCAL_OS_IMG}"
${DOCKER} tag "${KERNEL_IMG}" "${HTTPS_LOCAL_KERNEL_IMG}"

${DOCKER} login -u "${HTTP_USERNAME}" -p "${HTTP_PASSWORD}" "https://${HTTP_ADDR}"
${DOCKER} login -u "${HTTPS_USERNAME}" -p "${HTTPS_PASSWORD}" "https://${HTTPS_ADDR}"

# push in parallel, block until all finished
${DOCKER} push "${HTTP_LOCAL_OS_IMG}" &
${DOCKER} push "${HTTP_LOCAL_KERNEL_IMG}" &
${DOCKER} push "${HTTPS_LOCAL_OS_IMG}" &
${DOCKER} push "${HTTPS_LOCAL_KERNEL_IMG}" &
wait

${DOCKER} logout "http://${HTTP_ADDR}"
${DOCKER} logout "https://${HTTPS_ADDR}"

${DOCKER} rmi "${HTTP_LOCAL_OS_IMG}" "${HTTP_LOCAL_KERNEL_IMG}"
${DOCKER} rmi "${HTTPS_LOCAL_OS_IMG}" "${HTTPS_LOCAL_KERNEL_IMG}"
