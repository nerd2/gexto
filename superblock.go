package gexto

import (
	"os"
	"github.com/lunixbochs/struc"
	"io"
	"log"
)


type MoveExtent struct {
	Reserved    int32  `struc:"int32,little"`
	Donor_fd    int32  `struc:"int32,little"`
	Orig_start  uint64 `struc:"uint64,little"`
	Donor_start uint64 `struc:"uint64,little"`
	Len         uint64 `struc:"uint64,little"`
	Moved_len   uint64 `struc:"uint64,little"`
};

type ExtentHeader struct {
	Magic      int16 `struc:"int16,little"`
	Entries    int16 `struc:"int16,little"`
	Max        int16 `struc:"int16,little"`
	Depth      int16 `struc:"int16,little"`
	Generation int32 `struc:"int32,little"`
}

type Extent struct {
	Block    int32 `struc:"int32,little"`
	Len      int16 `struc:"int16,little"`
	Start_hi int16 `struc:"int16,little"`
	Start_lo int32 `struc:"int32,little"`
}

type DirectoryEntry2 struct {
	Inode int32 `struc:"int32,little"`
	Rec_len int16 `struc:"int16,little"`
	Name_len int8 `struc:"int8,sizeof=Name"`
	Flags int8 `struc:"int8"`
	Name string `struc:"[]byte"`
}

type Superblock struct {
	InodeCount         int32 `struc:"int32,little"`
	BlockCount_lo      int32 `struc:"int32,little"`
	R_blockCount_lo    int32 `struc:"int32,little"`
	Free_blockCount_lo int32 `struc:"int32,little"`
	Free_inodeCount    int32 `struc:"int32,little"`
	First_data_block   int32 `struc:"int32,little"`
	Log_block_size     int32 `struc:"int32,little"`
	Log_cluster_size   int32 `struc:"int32,little"`
	BlockPer_group     int32 `struc:"int32,little"`
	ClusterPer_group   int32 `struc:"int32,little"`
	InodePer_group     int32 `struc:"int32,little"`
	Mtime              int32 `struc:"int32,little"`
	Wtime              int32 `struc:"int32,little"`
	Mnt_count          int16 `struc:"int16,little"`
	Max_mnt_count      int16 `struc:"int16,little"`
	Magic              int16 `struc:"int16,little"`
	State              int16 `struc:"int16,little"`
	Errors             int16 `struc:"int16,little"`
	Minor_rev_level    int16 `struc:"int16,little"`
	Lastcheck          int32 `struc:"int32,little"`
	Checkinterval      int32 `struc:"int32,little"`
	Creator_os         int32 `struc:"int32,little"`
	Rev_level          int32 `struc:"int32,little"`
	Def_resuid         int16 `struc:"int16,little"`
	Def_resgid         int16 `struc:"int16,little"`
	// Dynamic_rev superblocks only
	First_ino              int32    `struc:"int32,little"`
	Inode_size             int16    `struc:"int16,little"`
	Block_group_nr         int16    `struc:"int16,little"`
	Feature_compat         int32    `struc:"int32,little"`
	Feature_incompat       int32    `struc:"int32,little"`
	Feature_ro_compat      int32    `struc:"int32,little"`
	Uuid                   [16]byte `struc:"byte"`
	Volume_name            [16]byte `struc:"byte"`
	Last_mounted           [64]byte `struc:"byte"`
	Algorithm_usage_bitmap int32    `struc:"int32,little"`
	// Performance hints
	Prealloc_blocks     byte  `struc:"byte"`
	Prealloc_dir_blocks byte  `struc:"byte"`
	Reserved_gdt_blocks int16 `struc:"int16,little"`
	// Journal

	Journal_Uuid       [16]byte  `struc:"byte"`
	Journal_inum       int32     `struc:"int32,little"`
	Journal_dev        int32     `struc:"int32,little"`
	Last_orphan        int32     `struc:"int32,little"`
	Hash_seed          [4]int32  `struc:"[4]int32,little"`
	Def_hash_version   byte      `struc:"byte"`
	Jnl_backup_type    byte      `struc:"byte"`
	Desc_size          int16     `struc:"int16,little"`
	Default_mount_opts int32     `struc:"int32,little"`
	First_meta_bg      int32     `struc:"int32,little"`
	MkfTime            int32     `struc:"int32,little"`
	Jnl_blocks         [17]int32 `struc:"[17]int32,little"`

	BlockCount_hi         int32     `struc:"int32,little"`
	R_blockCount_hi       int32     `struc:"int32,little"`
	Free_blockCount_hi    int32     `struc:"int32,little"`
	Min_extra_isize       int16     `struc:"int16,little"`
	Want_extra_isize      int16     `struc:"int16,little"`
	Flags                 int32     `struc:"int32,little"`
	Raid_stride           int16     `struc:"int16,little"`
	Mmp_update_interval   int16     `struc:"int16,little"`
	Mmp_block             int64     `struc:"int64,little"`
	Raid_stripe_width     int32     `struc:"int32,little"`
	Log_groupPer_flex     byte      `struc:"byte"`
	Checksum_type         byte      `struc:"byte"`
	Encryption_level      byte      `struc:"byte"`
	Reserved_pad          byte      `struc:"byte"`
	KbyteWritten          int64     `struc:"int64,little"`
	Snapshot_inum         int32     `struc:"int32,little"`
	Snapshot_id           int32     `struc:"int32,little"`
	Snapshot_r_blockCount int64     `struc:"int64,little"`
	Snapshot_list         int32     `struc:"int32,little"`
	Error_count           int32     `struc:"int32,little"`
	First_error_time      int32     `struc:"int32,little"`
	First_error_ino       int32     `struc:"int32,little"`
	First_error_block     int64     `struc:"int64,little"`
	First_error_func      [32]byte  `struc:"pad"`
	First_error_line      int32     `struc:"int32,little"`
	Last_error_time       int32     `struc:"int32,little"`
	Last_error_ino        int32     `struc:"int32,little"`
	Last_error_line       int32     `struc:"int32,little"`
	Last_error_block      int64     `struc:"int64,little"`
	Last_error_func       [32]byte  `struc:"pad"`
	Mount_opts            [64]byte  `struc:"pad"`
	Usr_quota_inum        int32     `struc:"int32,little"`
	Grp_quota_inum        int32     `struc:"int32,little"`
	Overhead_clusters     int32     `struc:"int32,little"`
	Backup_bgs            [2]int32  `struc:"[2]int32,little"`
	Encrypt_algos         [4]byte   `struc:"pad"`
	Encrypt_pw_salt       [16]byte  `struc:"pad"`
	Lpf_ino               int32     `struc:"int32,little"`
	Prj_quota_inum        int32     `struc:"int32,little"`
	Checksum_seed         int32     `struc:"int32,little"`
	Reserved              [98]int32 `struc:"[98]int32,little"`
	Checksum              int32     `struc:"int32,little"`
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

func getBlockGroupDescriptor(blockGroupNum int64,  sb *Superblock, dev *os.File) *GroupDescriptor {
	blockSize := sb.GetBlockSize()
	bgdtLocation := 1024/blockSize + 1

	log.Println(bgdtLocation,blockGroupNum,blockSize)
	bgd := &GroupDescriptor{}
	dev.Seek((bgdtLocation+blockGroupNum)*blockSize, 0)
	if sb.FeatureIncompat64bit() {
		struc.Unpack(dev, &bgd)
	} else {
		struc.Unpack(io.LimitReader(dev, 32), &bgd)
	}
	log.Printf("Read block group %d, contents:\n%+v\n", blockGroupNum, bgd)
	return bgd
}
