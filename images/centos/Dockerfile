FROM centos:7

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
