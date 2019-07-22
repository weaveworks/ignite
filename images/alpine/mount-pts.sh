#!/sbin/openrc-run

mkdir -p /dev/pts
mount devpts /dev/pts -t devpts || exit 0
