package constants

var BinaryDependencies = [...]string{
	"mount",
	"umount",
	"tar",
	"mkfs.ext4",
	"e2fsck",
	"resize2fs",
	"strings",
	"dmsetup",
	"ssh",
	"git",
}

var PathDependencies = [...]string{
	"/dev/mapper/control",
	"/dev/net/tun",
	"/dev/kvm",
}

var CNIDependencies = [...]string{
	"/opt/cni/bin/loopback",
	"/opt/cni/bin/bridge",
}
