package constants

const (
	// Ignite data base directory
	DATA_DIR = "/var/lib/firecracker"

	// Path to VM data directory containing a directory for each VM
	VM_DIR = DATA_DIR + "/vm"

	// Filename for the decompressed VM filesystem contents archive
	VM_FS_TAR = "fs.tar"

	// Filename for the filesystem disk image of a VM
	VM_FS_IMAGE = "fs.ext4"
)
