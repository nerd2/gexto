package gexto

type GroupDescriptor struct {
	Block_bitmap_lo      uint32 `struc:"uint32,little"`
	Inode_bitmap_lo      uint32 `struc:"uint32,little"`
	Inode_table_lo       uint32 `struc:"uint32,little"`
	Free_blocks_count_lo uint16 `struc:"uint16,little"`
	Free_inodes_count_lo uint16 `struc:"uint16,little"`
	Used_dirs_count_lo   uint16 `struc:"uint16,little"`
	Flags                uint16 `struc:"uint16,little"`
	Exclude_bitmap_lo    uint32 `struc:"uint32,little"`
	Block_bitmap_csum_lo uint16 `struc:"uint16,little"`
	Inode_bitmap_csum_lo uint16 `struc:"uint16,little"`
	Itable_unused_lo     uint16 `struc:"uint16,little"`
	Checksum             uint16 `struc:"uint16,little"`
	Block_bitmap_hi      uint32 `struc:"uint32,little"`
	Inode_bitmap_hi      uint32 `struc:"uint32,little"`
	Inode_table_hi       uint32 `struc:"uint32,little"`
	Free_blocks_count_hi uint16 `struc:"uint16,little"`
	Free_inodes_count_hi uint16 `struc:"uint16,little"`
	Used_dirs_count_hi   uint16 `struc:"uint16,little"`
	Itable_unused_hi     uint16 `struc:"uint16,little"`
	Exclude_bitmap_hi    uint32 `struc:"uint32,little"`
	Block_bitmap_csum_hi uint16 `struc:"uint16,little"`
	Inode_bitmap_csum_hi uint16 `struc:"uint16,little"`
	Reserved             uint32 `struc:"uint32,little"`
};

func (bgd *GroupDescriptor) GetInodeTableLoc(sb *Superblock) int64 {
	if sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_table_hi) << 32) | int64(bgd.Inode_table_lo)
	} else {
		return int64(bgd.Inode_table_lo)
	}
}
