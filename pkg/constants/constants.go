package constants

const (
	// Common Ignite prefix
	IGNITE_PREFIX = "ignite-"

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

	// Filename for VM overlay storage files
	OVERLAY_FILE = "overlay.dm"

	// In-container path for the mounted VM root overlay
	ROOT_DEV = "/ignite/rootfs"

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

	// Timeout in seconds to wait for VM shutdown before SIGKILL
	STOP_TIMEOUT = 20

	// Additional timeout in seconds for docker to wait for ignite to save and quit
	IGNITE_TIMEOUT = 10

	// In-container path for the firecracker socket
	SOCKET_PATH = "/tmp/firecracker.sock"

	// Common VM kernel parameters
	VM_KERNEL_ARGS = "console=ttyS0 reboot=k panic=1 pci=off"

	// In-container path for the firecracker log FIFO
	LOG_FIFO = "/tmp/firecracker_log.fifo"

	// In-container path for the firecracker metrics FIFO
	METRICS_FIFO = "/tmp/firecracker_metrics.fifo"

	// Log level for the firecracker VM
	VM_LOG_LEVEL = "Error"
)
