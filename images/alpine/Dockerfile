FROM scratch
# TODO: Base this on the alpine docker image, not the FC ext4 image
ADD alpine.tar /
# Add an SSH server and start it automatically
RUN apk add \
        openssh \
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
