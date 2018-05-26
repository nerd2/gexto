package gexto

type GroupDescriptor struct {
	Block_bitmap_lo      int32 `struc:"int32,little"`
	Inode_bitmap_lo      int32 `struc:"int32,little"`
	Inode_table_lo       int32 `struc:"int32,little"`
	Free_blocks_count_lo int16 `struc:"int16,little"`
	Free_inodes_count_lo int16 `struc:"int16,little"`
	Used_dirs_count_lo   int16 `struc:"int16,little"`
	Flags                int16 `struc:"int16,little"`
	Exclude_bitmap_lo    int32 `struc:"int32,little"`
	Block_bitmap_csum_lo int16 `struc:"int16,little"`
	Inode_bitmap_csum_lo int16 `struc:"int16,little"`
	Itable_unused_lo     int16 `struc:"int16,little"`
	Checksum             int16 `struc:"int16,little"`
	Block_bitmap_hi      int32 `struc:"int32,little"`
	Inode_bitmap_hi      int32 `struc:"int32,little"`
	Inode_table_hi       int32 `struc:"int32,little"`
	Free_blocks_count_hi int16 `struc:"int16,little"`
	Free_inodes_count_hi int16 `struc:"int16,little"`
	Used_dirs_count_hi   int16 `struc:"int16,little"`
	Itable_unused_hi     int16 `struc:"int16,little"`
	Exclude_bitmap_hi    int32 `struc:"int32,little"`
	Block_bitmap_csum_hi int16 `struc:"int16,little"`
	Inode_bitmap_csum_hi int16 `struc:"int16,little"`
	Reserved             int32 `struc:"int32,little"`
};

func (bgd *GroupDescriptor) GetInodeTableLoc(sb *Superblock) int64 {
	if sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_table_hi) << 32) | int64(bgd.Inode_table_lo)
	} else {
		return int64(bgd.Inode_table_lo)
	}
}
