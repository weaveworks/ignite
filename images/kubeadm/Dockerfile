FROM weaveworks/ignite-ubuntu:latest
# Install dependencies
RUN apt-get update && apt-get install -y \
        apt-transport-https \
        curl \
        docker.io \
        gnupg2 \
        jq \
    && apt-get clean

# Install gcloud
RUN curl -sSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb http://packages.cloud.google.com/apt cloud-sdk-bionic main" > /etc/apt/sources.list.d/google-cloud-sdk.list && \
    apt-get update && apt-get install -y google-cloud-sdk

# Install k8s
ENV KUBERNETES_VERSION=v1.15.0-rc.1
COPY install.sh /
RUN /install.sh install ci/latest-1.15 ${KUBERNETES_VERSION}
RUN apt-get update && apt-get install -f -y

# Remove gcloud
RUN apt-get purge -y google-cloud-sdk
