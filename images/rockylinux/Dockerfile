ARG RELEASE=8

FROM rockylinux/rockylinux:${RELEASE}

# Shadow the bogus /etc/resolv.conf of centos:8 by copying a blank file over it
COPY resolv.conf /etc/

# Install common utilities
RUN dnf -y install --setopt=install_weak_deps=False --setopt=tsflags=nodocs \
        iproute \
        iputils \
        openssh-server \
        net-tools \
        procps-ng \
        sudo \
        wget && \
    dnf clean all

# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo "root:root" | chpasswd
