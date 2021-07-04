package constants

import "time"

const (
	// Path to directory containing a subdirectory for each VM
	VM_DIR = DATA_DIR + "/vm"

	// Path where ignited stores its manifests
	MANIFEST_DIR = "/etc/firecracker/manifests"

	// Default values for VM options
	VM_DEFAULT_CPUS        = 1
	VM_DEFAULT_MEMORY      = 512 * MB
	VM_DEFAULT_SIZE        = 4 * GB
	VM_DEFAULT_KERNEL_ARGS = "console=ttyS0 reboot=k panic=1 pci=off ip=dhcp"

	// SSH key template for VMs
	VM_SSH_KEY_TEMPLATE = "id_%s"

	// TODO: remove this when the old dm code is removed
	OVERLAY_FILE = "overlay.dm"

	// Prometheus socket filename
	PROMETHEUS_SOCKET = "prometheus.sock"

	// Where the VM specification is located inside of the container
	IGNITE_SPAWN_VM_FILE_PATH = "/vm.json"

	// Where the vmlinux kernel is located inside of the container
	IGNITE_SPAWN_VMLINUX_FILE_PATH = "/vmlinux"

	// Subdirectory for volumes to be forwarded into the VM
	IGNITE_SPAWN_VOLUME_DIR = "/volumes"

	// DEFAULT_SANDBOX_IMAGE_NAME is the name of the default sandbox container
	// image to be used.
	DEFAULT_SANDBOX_IMAGE_NAME = "weaveworks/ignite"

	// DEFAULT_SANDBOX_IMAGE_NAME is the name of the default sandbox container
	// image to be used.
	DEFAULT_SANDBOX_IMAGE_TAG = "dev"

	// IGNITE_INTERFACE_ANNOTATION is the annotation prefix to store a list of extra interfaces
	IGNITE_INTERFACE_ANNOTATION = "ignite.weave.works/interface/"

	// IGNITE_SANDBOX_ENV_VAR is the annotation prefix to store a list of env variables
	IGNITE_SANDBOX_ENV_VAR = "ignite.weave.works/sanbox-env/"

	// IGNITE_SPAWN_TIMEOUT determines how long to wait for spawn to start up
	IGNITE_SPAWN_TIMEOUT = 2 * time.Minute
)
