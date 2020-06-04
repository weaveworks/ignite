ARG RELEASE

FROM centos:${RELEASE}

# Shadow the bogus /etc/resolv.conf of centos:8 by copying a blank file over it
COPY resolv.conf /etc/

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
RUN echo "root:root" | chpasswd
