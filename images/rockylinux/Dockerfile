ARG DIGEST
FROM rockylinux/rockylinux@sha256:${DIGEST}

# If we're building for another architecture than amd64, this let's us emulate an other platform's docker build.
# If we're building normally, for amd64, this line is removed
COPY qemu-QEMUARCH-static /usr/bin/

# Shadow the bogus /etc/resolv.conf of centos:8 by copying a blank file over it
COPY resolv.conf /etc/

# Install common utilities
RUN dnf -y install --setopt=install_weak_deps=False --setopt=tsflags=nodocs \
        iproute \
        iputils \
        openssh-server \
        net-tools \
        procps-ng \
        wget && \
    dnf clean all

# Create the following files, but unset them
RUN echo "" > /etc/machine-id && echo "" > /var/lib/dbus/machine-id

# This container image doesn't have locales installed. Disable forwarding the
# user locale env variables or we get warnings such as:
#  bash: warning: setlocale: LC_ALL: cannot change locale
RUN sed -i -e 's/^AcceptEnv LANG LC_\*$/#AcceptEnv LANG LC_*/' /etc/ssh/sshd_config

# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo "root:root" | chpasswd
