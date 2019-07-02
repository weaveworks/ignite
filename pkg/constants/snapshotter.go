package constants

const (
	// Default maximum size for the pool (if using a physical device, it's size will be used instead)
	POOL_DATA_SIZE_BYTES = 100 * GB // 100 GB

	// Default allocation size for the pool, should be between
	// 128 (64KB) and 2097152 (1GB). 128 is recommended if
	// snapshotting a lot (like we do with layers).
	POOL_ALLOCATION_SIZE_SECTORS = 128

	// Additional space for volumes to accommodate the ext4 partition
	POOL_VOLUME_EXTRA_SIZE = 100 * MB

	// Base directory for snapshotter data
	SNAPSHOTTER_DIR = DATA_DIR + "/snapshotter"

	// Paths to the default data and metadata backing files
	SNAPSHOTTER_METADATA_PATH = SNAPSHOTTER_DIR + "/metadata.dm"
	SNAPSHOTTER_DATA_PATH     = SNAPSHOTTER_DIR + "/data.dm"
)
