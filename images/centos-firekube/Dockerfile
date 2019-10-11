ARG RELEASE

FROM centos:${RELEASE}

# Install common utilities
RUN yum -y install \
        iproute \
        iputils \
        openssh-server \
        net-tools \
        procps-ng \
        sudo \
        wget

ADD  https://github.com/chanwit/containerd/releases/download/v1.2.9-p1/ctr /usr/bin/ctr129
RUN  chmod +x /usr/bin/ctr129
COPY kubernetes.repo /etc/yum.repos.d/kubernetes.repo
COPY docker-ce.repo /etc/yum.repos.d/docker-ce.repo
COPY preload_images.sh /tmp/images/preload_images.sh
RUN  chmod +x /tmp/images/preload_images.sh
COPY preloaded_images /tmp/images/preloaded_images
COPY tmp/*.tar /tmp/images/

RUN yum -y install \
    yum-plugin-downloadonly \
	yum-plugin-versionlock

RUN yum install --downloadonly device-mapper-persistent-data lvm2
RUN yum install --downloadonly docker-ce-18.09.7
RUN yum install --downloadonly kubelet-1.14.1 --disableexcludes kubernetes
RUN yum install --downloadonly kubeadm-1.14.1 --disableexcludes kubernetes

RUN yum -y install containerd.io libseccomp

WORKDIR /tmp/images
RUN /usr/bin/containerd & \
    sleep 2 && \
    ctr129 version && \
    bash -x preload_images.sh

RUN rm -f /etc/docker/daemon.json
RUN rm -rf /tmp/images/
RUN rm -f /usr/bin/ctr129

FROM centos:${RELEASE}

# Install common utilities
RUN yum -y install \
        iproute \
        iputils \
        openssh-server \
        net-tools \
        procps-ng \
        sudo \
        wget && \
    yum clean all

# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo root | passwd --stdin root

COPY --from=0 /var/cache/yum /var/cache/yum
COPY --from=0 /var/lib/containerd /var/lib/containerd
# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo root | passwd --stdin root