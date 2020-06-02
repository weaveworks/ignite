# Ubuntu 20.04 was also tested, but didn't perform very well (sshd took a long time to start), so we're sticking with Ubuntu 18.04 still
FROM weaveworks/ignite-ubuntu:18.04
# Install dependencies. Use containerd for running the containers (for better performance)
RUN apt-get update && apt-get install -y --no-install-recommends \
        apt-transport-https \
        containerd \
        curl \
        gnupg2 \
        jq \
    && apt-get clean

# Install k8s v1.18.1 locally
ENV KUBERNETES_VERSION=v1.18.3
COPY install.sh /
RUN /install.sh install release/stable-1.18 ${KUBERNETES_VERSION}
# Docker sets this automatically, but not containerd.
# It is required when running kubeadm.
RUN echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf
