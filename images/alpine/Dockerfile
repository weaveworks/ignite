FROM scratch
# Bases on the alpine edge tar file
ADD alpine.tar /
# Add an SSH server and start it automatically
RUN apk add \
        openssh \
        openrc \
        util-linux \
        udev \
        bash

COPY mount-pts.sh /etc/init.d/devpts
RUN rc-update add sshd && \
    rc-update add udev && \
    rc-update add sysfs && \
    rc-update add devpts

RUN echo "exit 0" > /etc/init.d/networking
RUN echo "PermitTTY yes" >> /etc/ssh/sshd_config
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
# set the default root password for agetty
RUN echo root:root | chpasswd
# configure agetty
RUN ln -s /etc/init.d/agetty /etc/init.d/agetty.ttyS0
RUN ln -s /etc/init.d/agetty.ttyS0 /etc/runlevels/default
RUN ln -s /etc/init.d/agetty.ttyS0 /etc/runlevels/nonetwork
