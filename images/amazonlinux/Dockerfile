ARG RELEASE

FROM amazonlinux:${RELEASE}

# Install common utilities
RUN yum -y install \
        hostname \
        iproute \
        iputils \
        net-tools \
        openssh-server \
        procps-ng \
        sudo \
        systemd \
        wget \
    && yum clean all

# TODO: Set the root password to root when logging in through the VM's ttyS0 console
