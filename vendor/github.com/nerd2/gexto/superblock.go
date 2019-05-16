package gexto

import (
	"github.com/lunixbochs/struc"
)

type Superblock struct {
	InodeCount         uint32 `struc:"uint32,little"`
	BlockCount_lo      uint32 `struc:"uint32,little"`
	R_blockCount_lo    uint32 `struc:"uint32,little"`
	Free_blockCount_lo uint32 `struc:"uint32,little"`
	Free_inodeCount    uint32 `struc:"uint32,little"`
	First_data_block   uint32 `struc:"uint32,little"`
	Log_block_size     uint32 `struc:"uint32,little"`
	Log_cluster_size   uint32 `struc:"uint32,little"`
	BlockPer_group     uint32 `struc:"uint32,little"`
	ClusterPer_group   uint32 `struc:"uint32,little"`
	InodePer_group     uint32 `struc:"uint32,little"`
	Mtime              uint32 `struc:"uint32,little"`
	Wtime              uint32 `struc:"uint32,little"`
	Mnt_count          uint16 `struc:"uint16,little"`
	Max_mnt_count      uint16 `struc:"uint16,little"`
	Magic              uint16 `struc:"uint16,little"`
	State              uint16 `struc:"uint16,little"`
	Errors             uint16 `struc:"uint16,little"`
	Minor_rev_level    uint16 `struc:"uint16,little"`
	Lastcheck          uint32 `struc:"uint32,little"`
	Checkinterval      uint32 `struc:"uint32,little"`
	Creator_os         uint32 `struc:"uint32,little"`
	Rev_level          uint32 `struc:"uint32,little"`
	Def_resuid         uint16 `struc:"uint16,little"`
	Def_resgid         uint16 `struc:"uint16,little"`
	// Dynamic_rev superblocks only
	First_ino              uint32    `struc:"uint32,little"`
	Inode_size             uint16    `struc:"uint16,little"`
	Block_group_nr         uint16    `struc:"uint16,little"`
	Feature_compat         uint32    `struc:"uint32,little"`
	Feature_incompat       uint32    `struc:"uint32,little"`
	Feature_ro_compat      uint32    `struc:"uint32,little"`
	Uuid                   [16]byte `struc:"[16]byte"`
	Volume_name            [16]byte `struc:"[16]byte"`
	Last_mounted           [64]byte `struc:"[64]byte"`
	Algorithm_usage_bitmap uint32    `struc:"uint32,little"`
	// Performance hints
	Prealloc_blocks     byte  `struc:"byte"`
	Prealloc_dir_blocks byte  `struc:"byte"`
	Reserved_gdt_blocks uint16 `struc:"uint16,little"`
	// Journal

	Journal_Uuid       [16]byte  `struc:"[16]byte"`
	Journal_inum       uint32     `struc:"uint32,little"`
	Journal_dev        uint32     `struc:"uint32,little"`
	Last_orphan        uint32     `struc:"uint32,little"`
	Hash_seed          [4]uint32  `struc:"[4]uint32,little"`
	Def_hash_version   byte      `struc:"byte"`
	Jnl_backup_type    byte      `struc:"byte"`
	Desc_size          uint16     `struc:"uint16,little"`
	Default_mount_opts uint32     `struc:"uint32,little"`
	First_meta_bg      uint32     `struc:"uint32,little"`
	MkfTime            uint32     `struc:"uint32,little"`
	Jnl_blocks         [17]uint32 `struc:"[17]uint32,little"`

	BlockCount_hi         uint32     `struc:"uint32,little"`
	R_blockCount_hi       uint32     `struc:"uint32,little"`
	Free_blockCount_hi    uint32     `struc:"uint32,little"`
	Min_extra_isize       uint16     `struc:"uint16,little"`
	Want_extra_isize      uint16     `struc:"uint16,little"`
	Flags                 uint32     `struc:"uint32,little"`
	Raid_stride           uint16     `struc:"uint16,little"`
	Mmp_update_interval   uint16     `struc:"uint16,little"`
	Mmp_block             uint64     `struc:"uint64,little"`
	Raid_stripe_width     uint32     `struc:"uint32,little"`
	Log_groupPer_flex     byte      `struc:"byte"`
	Checksum_type         byte      `struc:"byte"`
	Encryption_level      byte      `struc:"byte"`
	Reserved_pad          byte      `struc:"byte"`
	KbyteWritten          uint64     `struc:"uint64,little"`
	Snapshot_inum         uint32     `struc:"uint32,little"`
	Snapshot_id           uint32     `struc:"uint32,little"`
	Snapshot_r_blockCount uint64     `struc:"uint64,little"`
	Snapshot_list         uint32     `struc:"uint32,little"`
	Error_count           uint32     `struc:"uint32,little"`
	First_error_time      uint32     `struc:"uint32,little"`
	First_error_ino       uint32     `struc:"uint32,little"`
	First_error_block     uint64     `struc:"uint64,little"`
	First_error_func      [32]byte  `struc:"[32]pad"`
	First_error_line      uint32     `struc:"uint32,little"`
	Last_error_time       uint32     `struc:"uint32,little"`
	Last_error_ino        uint32     `struc:"uint32,little"`
	Last_error_line       uint32     `struc:"uint32,little"`
	Last_error_block      uint64     `struc:"uint64,little"`
	Last_error_func       [32]byte  `struc:"[32]pad"`
	Mount_opts            [64]byte  `struc:"[64]pad"`
	Usr_quota_inum        uint32     `struc:"uint32,little"`
	Grp_quota_inum        uint32     `struc:"uint32,little"`
	Overhead_clusters     uint32     `struc:"uint32,little"`
	Backup_bgs            [2]uint32  `struc:"[2]uint32,little"`
	Encrypt_algos         [4]byte   `struc:"[4]pad"`
	Encrypt_pw_salt       [16]byte  `struc:"[16]pad"`
	Lpf_ino               uint32     `struc:"uint32,little"`
	Prj_quota_inum        uint32     `struc:"uint32,little"`
	Checksum_seed         uint32     `struc:"uint32,little"`
	Reserved              [98]uint32 `struc:"[98]uint32,little"`
	Checksum              uint32     `struc:"uint32,little"`
	address               int64
	fs                    *fs
	numBlockGroups        int64
};

func (sb *Superblock) FeatureCompatDir_prealloc() bool  { return (sb.Feature_compat&FEATURE_COMPAT_DIR_PREALLOC != 0) }
func (sb *Superblock) FeatureCompatImagic_inodes() bool { return (sb.Feature_compat&FEATURE_COMPAT_IMAGIC_INODES != 0) }
func (sb *Superblock) FeatureCompatHas_journal() bool   { return (sb.Feature_compat&FEATURE_COMPAT_HAS_JOURNAL != 0) }
func (sb *Superblock) FeatureCompatExt_attr() bool      { return (sb.Feature_compat&FEATURE_COMPAT_EXT_ATTR != 0) }
func (sb *Superblock) FeatureCompatResize_inode() bool  { return (sb.Feature_compat&FEATURE_COMPAT_RESIZE_INODE != 0) }
func (sb *Superblock) FeatureCompatDir_index() bool     { return (sb.Feature_compat&FEATURE_COMPAT_DIR_INDEX != 0) }
func (sb *Superblock) FeatureCompatSparse_super2() bool { return (sb.Feature_compat&FEATURE_COMPAT_SPARSE_SUPER2 != 0) }

func (sb *Superblock) FeatureRoCompatSparse_super() bool  { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_SPARSE_SUPER != 0) }
func (sb *Superblock) FeatureRoCompatLarge_file() bool    { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_LARGE_FILE != 0) }
func (sb *Superblock) FeatureRoCompatBtree_dir() bool     { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_BTREE_DIR != 0) }
func (sb *Superblock) FeatureRoCompatHuge_file() bool     { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_HUGE_FILE != 0) }
func (sb *Superblock) FeatureRoCompatGdt_csum() bool      { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_GDT_CSUM != 0) }
func (sb *Superblock) FeatureRoCompatDir_nlink() bool     { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_DIR_NLINK != 0) }
func (sb *Superblock) FeatureRoCompatExtra_isize() bool   { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_EXTRA_ISIZE != 0) }
func (sb *Superblock) FeatureRoCompatQuota() bool         { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_QUOTA != 0) }
func (sb *Superblock) FeatureRoCompatBigalloc() bool      { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_BIGALLOC != 0) }
func (sb *Superblock) FeatureRoCompatMetadata_csum() bool { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_METADATA_CSUM != 0) }
func (sb *Superblock) FeatureRoCompatReadonly() bool      { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_READONLY != 0) }
func (sb *Superblock) FeatureRoCompatProject() bool       { return (sb.Feature_ro_compat&FEATURE_RO_COMPAT_PROJECT != 0) }

func (sb *Superblock) FeatureIncompat64bit() bool       { return (sb.Feature_incompat&FEATURE_INCOMPAT_64BIT != 0) }
func (sb *Superblock) FeatureIncompatCompression() bool { return (sb.Feature_incompat&FEATURE_INCOMPAT_COMPRESSION != 0) }
func (sb *Superblock) FeatureIncompatFiletype() bool    { return (sb.Feature_incompat&FEATURE_INCOMPAT_FILETYPE != 0) }
func (sb *Superblock) FeatureIncompatRecover() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_RECOVER != 0) }
func (sb *Superblock) FeatureIncompatJournal_dev() bool { return (sb.Feature_incompat&FEATURE_INCOMPAT_JOURNAL_DEV != 0) }
func (sb *Superblock) FeatureIncompatMeta_bg() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_META_BG != 0) }
func (sb *Superblock) FeatureIncompatExtents() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_EXTENTS != 0) }
func (sb *Superblock) FeatureIncompatMmp() bool         { return (sb.Feature_incompat&FEATURE_INCOMPAT_MMP != 0) }
func (sb *Superblock) FeatureIncompatFlex_bg() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_FLEX_BG != 0) }
func (sb *Superblock) FeatureIncompatEa_inode() bool    { return (sb.Feature_incompat&FEATURE_INCOMPAT_EA_INODE != 0) }
func (sb *Superblock) FeatureIncompatDirdata() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_DIRDATA != 0) }
func (sb *Superblock) FeatureIncompatCsum_seed() bool   { return (sb.Feature_incompat&FEATURE_INCOMPAT_CSUM_SEED != 0) }
func (sb *Superblock) FeatureIncompatLargedir() bool    { return (sb.Feature_incompat&FEATURE_INCOMPAT_LARGEDIR != 0) }
func (sb *Superblock) FeatureIncompatInline_data() bool { return (sb.Feature_incompat&FEATURE_INCOMPAT_INLINE_DATA != 0) }
func (sb *Superblock) FeatureIncompatEncrypt() bool     { return (sb.Feature_incompat&FEATURE_INCOMPAT_ENCRYPT != 0) }

func (sb *Superblock) GetBlockCount() int64 {
	if sb.FeatureIncompat64bit() {
		return (int64(sb.BlockCount_hi) << 32) | int64(sb.BlockCount_lo)
	} else {
		return int64(sb.BlockCount_lo)
	}
}

func (sb *Superblock) GetBlockSize() int64 {
	return int64(1024 << uint(sb.Log_block_size))
}

func (sb *Superblock) UpdateCsumAndWriteback() {
	cs := NewChecksummer(sb)

	size, _ := struc.Sizeof(sb)
	struc.Pack(LimitWriter(cs, int64(size) - 4), sb)
	sb.Checksum = cs.Get()

	sb.fs.dev.Seek(sb.address, 0)
	struc.Pack(sb.fs.dev, sb)
}

func (sb *Superblock) GetGroupsPerFlex() int64 {
	return 1 << sb.Log_groupPer_flex
}