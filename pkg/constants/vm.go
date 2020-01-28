package constants

const (
	// Path to directory containing a subdirectory for each VM
	VM_DIR = DATA_DIR + "/vm"

	// Path where ignited stores its manifests
	MANIFEST_DIR = "/etc/firecracker/manifests"

	// Default values for VM options
	VM_DEFAULT_CPUS   = 1
	VM_DEFAULT_MEMORY = 512 * MB
	VM_DEFAULT_SIZE   = 4 * GB
	// Refer to https://github.com/firecracker-microvm/firecracker/blob/master/src/vmm/src/vmm_config/boot_source.rs
	// TODO: The Firecracker team don't use console=ttyS0 anymore, but 8250.nr_uarts=0 instead
	// See https://github.com/firecracker-microvm/firecracker/pull/313 for more info
	VM_DEFAULT_KERNEL_ARGS = "console=ttyS0 reboot=k panic=1 pci=off ip=dhcp i8042.noaux i8042.nomux i8042.nopnp i8042.dumbkbd"

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
)
