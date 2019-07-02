package constants

const (
	// Common Ignite prefix
	IGNITE_PREFIX = "ignite-"

	// Ignite data base directory
	DATA_DIR = "/var/lib/firecracker"

	// Permissions for the data directory and its subdirectories
	DATA_DIR_PERM = 0755

	// Filename for metadata files
	METADATA = "metadata.json"

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

	// In-container path for the firecracker log FIFO
	LOG_FIFO = "/tmp/firecracker_log.fifo"

	// In-container path for the firecracker metrics FIFO
	METRICS_FIFO = "/tmp/firecracker_metrics.fifo"

	// Log level for the firecracker VM
	VM_LOG_LEVEL = "Error"
)
