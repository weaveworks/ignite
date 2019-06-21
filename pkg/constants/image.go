package constants

const (
	// Path to directory containing a subdirectory for each image
	IMAGE_DIR = DATA_DIR + "/image"

	// Filename for the image file containing the image filesystem
	IMAGE_FS = "image.ext4"

	// Filename for the thin provisioning metadata file
	IMAGE_THINMETADATA = "metadata.dm"

	// Filename for the thin provisioning data file
	IMAGE_THINDATA = "data.dm"

	// The default kernel for VMs
	DEFAULT_KERNEL = "weaveworks/ignite-kernel:4.19.47"
)
