package constants

const (
	// Path to directory containing a subdirectory for each VM
	VM_DIR = DATA_DIR + "/vm"

	// Default values for VM options
	VM_DEFAULT_CPUS        = 1
	VM_DEFAULT_MEMORY      = 512 * MB
	VM_DEFAULT_SIZE        = 4 * GB
	VM_DEFAULT_KERNEL_ARGS = "console=ttyS0 reboot=k panic=1 pci=off ip=dhcp"

	// SSH key template for VMs
	VM_SSH_KEY_TEMPLATE = "id_%s"

	// TODO: remove this when the old dm code is removed
	OVERLAY_FILE = "overlay.dm"
)
