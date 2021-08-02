ARG RELEASE

FROM opensuse/${RELEASE}

# Install common utilities
RUN zypper -n install \
        iproute \
        iputils \
        openssh \
        net-tools \
        systemd-sysvinit \
        udev \
        sudo \
        shadow \
        wget && \
    zypper clean --all

# systemctl enable creates the symlinks to enable the service
# it doesn't need systemd running
RUN systemctl enable sshd

# Set the root password to root when logging in through the VM's ttyS0 console
RUN echo "root:root" | chpasswd

