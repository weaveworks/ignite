#!/usr/bin/env bash

set -e

# This script runs a local private docker registry with self-signed certificate
# and basic auth.

REGISTRY_SECRET_PATH=/tmp/ignite-test-registry
PRIVATE_KEY=${REGISTRY_SECRET_PATH}/certs/domain.key
CERT=${REGISTRY_SECRET_PATH}/certs/domain.crt
HTPASSWD=${REGISTRY_SECRET_PATH}/auth/htpasswd
USERNAME=testuser
PASSWORD=testpassword
REGISTRY_ADDRESS=https://localhost:5000
OS_IMG=weaveworks/ignite-ubuntu:latest
LOCAL_OS_IMG=localhost:5000/weaveworks/ignite-ubuntu:test
KERNEL_IMG=weaveworks/ignite-kernel:5.4.108
LOCAL_KERNEL_IMG=localhost:5000/weaveworks/ignite-kernel:test

# Clear any existing registry secret and create new directories.
rm -rf ${REGISTRY_SECRET_PATH}
mkdir -p ${REGISTRY_SECRET_PATH}/{certs,auth}

# Generate key and cert.
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=US/ST=Foo/L=Bar/O=Weave" \
    -keyout ${PRIVATE_KEY} -out ${CERT}
chmod 400 ${PRIVATE_KEY}

# Create htpasswd file.
docker run --rm \
  --entrypoint htpasswd \
  httpd:2 -Bbn ${USERNAME} ${PASSWORD} > ${HTPASSWD}

# Run the registry.
docker run -d --rm \
  --name registry \
  -v ${REGISTRY_SECRET_PATH}/auth:/auth \
  -v ${REGISTRY_SECRET_PATH}/certs:/certs \
  -e REGISTRY_AUTH=htpasswd \
  -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
  -p 5000:5000 \
  registry:2

# Login, push test images to download in tests and logout.
docker login -u ${USERNAME} -p ${PASSWORD} ${REGISTRY_ADDRESS}
docker pull ${OS_IMG}
docker pull ${KERNEL_IMG}
docker tag ${OS_IMG} ${LOCAL_OS_IMG}
docker tag ${KERNEL_IMG} ${LOCAL_KERNEL_IMG}
docker push ${LOCAL_OS_IMG}
docker push ${LOCAL_KERNEL_IMG}
docker rmi ${LOCAL_OS_IMG} ${LOCAL_KERNEL_IMG}
docker logout ${REGISTRY_ADDRESS}
