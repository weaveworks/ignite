FROM weaveworks/ignite-kernel:4.19.47 AS kernel
FROM centos:7

# Install common utilities
RUN yum -y install \
        sudo \
        procps-ng \
        net-tools \
        iproute \
        iputils \
        openssh-server \
        wget && \
    yum clean all

# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo root | passwd --stdin root

# Copy in the kernel overlay files
COPY --from=kernel / /
