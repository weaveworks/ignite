ARG GOARCH="amd64"

# Ubuntu 20.04 was also tested, but didn't perform very well (sshd took a long time to start), so we're sticking with Ubuntu 18.04 still
FROM weaveworks/ignite-ubuntu:18.04
ARG GOARCH="amd64"
ARG RELEASE
ARG BINARY_REF

# Install dependencies. Use containerd for running the containers (for better performance)
RUN apt-get update && apt-get install -y --no-install-recommends \
        apt-transport-https \
        containerd \
        curl \
        gnupg2 \
        jq \
    && apt-get clean

# Install k8s locally
COPY ./install.sh /
RUN /install.sh install "${BINARY_REF}" "${RELEASE}" "${GOARCH}"
# Docker sets this automatically, but not containerd.
# It is required when running kubeadm.
RUN echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf
