package constants

const (
	// Ignite data base directory
	DATA_DIR = "/var/lib/firecracker"

	// Path to data directory containing a directory for each image
	IMAGE_DIR = DATA_DIR + "/image"

	// Filename for the decompressed filesystem contents archive
	IMAGE_TAR = "image.tar"

	// Filename for the disk image containing the filesystem
	IMAGE_FS = "image.ext4"

	// Filename for metadata files
	METADATA = "metadata.json"

	// Directory for hosting VM instances
	VM_DIR = DATA_DIR + "/vm"

	// Directory containing VM kernels
	KERNEL_DIR = DATA_DIR + "/kernel"

	// Kernel filename
	KERNEL_FILE = "vmlinux"

	// DHCP infinite lease time
	DHCP_INFINITE_LEASE = "4294967295s"

	// TAP adapter prefix in the parent container
	TAP_PREFIX = "vm_"

	// Bridge device prefix in the parent container
	BRIDGE_PREFIX = "br_"
)
