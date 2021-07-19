package constants

const (
	// Path to directory containing a subdirectory for each kernel
	KERNEL_DIR = DATA_DIR + "/kernel"

	// Kernel filename
	KERNEL_FILE = "vmlinux"

	// Filename for the tar containing the kernel filesystem
	KERNEL_TAR = "kernel.tar"

	// The kernel image name to be used as the default
	DEFAULT_KERNEL_IMAGE_NAME = "weaveworks/ignite-kernel"

	// The kernel image tag to be used as the default
	DEFAULT_KERNEL_IMAGE_TAG = "5.10.51"
)
