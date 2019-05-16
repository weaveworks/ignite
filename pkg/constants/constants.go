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
)
