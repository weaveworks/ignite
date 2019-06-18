FROM weaveworks/ignite-ubuntu:latest
# Install dependencies
RUN apt-get update && apt-get install -y \
        apt-transport-https \
        curl \
        docker.io \
        gnupg2 \
        jq \
    && apt-get clean

# Install k8s v1.15.0-rc.1 locally
ENV KUBERNETES_VERSION=v1.15.0-rc.1
COPY install.sh /
RUN /install.sh install ci/latest-1.15 ${KUBERNETES_VERSION}
