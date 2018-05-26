package gexto

import (
	"github.com/lunixbochs/struc"
	"math/bits"
)

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
	fs                   *fs
	num                  int64
	address              int64
};

func (bgd *GroupDescriptor) GetInodeBitmapLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_bitmap_hi) << 32) | int64(bgd.Inode_bitmap_lo)
	} else {
		return int64(bgd.Inode_bitmap_lo)
	}
}

func (bgd *GroupDescriptor) GetInodeTableLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_table_hi) << 32) | int64(bgd.Inode_table_lo)
	} else {
		return int64(bgd.Inode_table_lo)
	}
}

func (bgd *GroupDescriptor) GetBlockBitmapLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Block_bitmap_hi) << 32) | int64(bgd.Block_bitmap_lo)
	} else {
		return int64(bgd.Block_bitmap_lo)
	}
}

func (bgd *GroupDescriptor) UpdateCsumAndWriteback() {
	cs := NewChecksummer(bgd.fs.sb)

	cs.Write(bgd.fs.sb.Uuid[:])
	cs.WriteUint32(uint32(bgd.num))
	bgd.Checksum = 0
	struc.Pack(cs, bgd)
	bgd.Checksum = uint16(cs.Get() & 0xFFFF)

	bgd.fs.dev.Seek(bgd.address, 0)
	struc.Pack(bgd.fs.dev, bgd)
}

func(bgd *GroupDescriptor) GetFreeInode() *Inode {
	// Find free inode in bitmap
	start := bgd.GetInodeBitmapLoc() * bgd.fs.sb.GetBlockSize()
	bgd.fs.dev.Seek(start, 0)
	subInodeNum := int64(-1)
	for i := 0; i < int(bgd.fs.sb.InodePer_group/8); i++ {
		b := make([]byte, 1)
		bgd.fs.dev.Read(b)
		if b[0] != 0xFF {
			bitNum := bits.TrailingZeros8(^b[0])
			subInodeNum = int64(i) * 8 + int64(bitNum)
			b[0] |= 1 << uint(bitNum)
			bgd.fs.dev.Seek(-1, 1)
			bgd.fs.dev.Write(b)
			break
		}
	}

	if subInodeNum < 0 {
		return nil
	}

	// Update inode bitmap checksum
	checksummer := NewChecksummer(bgd.fs.sb)
	checksummer.Write(bgd.fs.sb.Uuid[:])
	bgd.fs.dev.Seek(start, 0)
	b := make([]byte, int64(bgd.fs.sb.InodePer_group) / 8)
	bgd.fs.dev.Read(b)
	checksummer.Write(b)
	bgd.Inode_bitmap_csum_lo = uint16(checksummer.Get() & 0xFFFF)
	bgd.Inode_bitmap_csum_hi = uint16(checksummer.Get() >> 16)

	bgd.Free_inodes_count_lo--
	bgd.Itable_unused_lo--
	bgd.UpdateCsumAndWriteback()

	bgd.fs.sb.Free_inodeCount--
	bgd.fs.sb.UpdateCsumAndWriteback()

	// Insert in Inode table
	inode := &Inode{
		Mode: 0x41FF,
		Links_count: 1,
		Flags: 524288, //TODO: what
		BlockOrExtents: [60]byte{0x0a, 0xf3, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00},
		fs: bgd.fs,
		address: bgd.GetInodeTableLoc() * bgd.fs.sb.GetBlockSize() + subInodeNum * int64(bgd.fs.sb.Inode_size),
		num: 1 + bgd.num * int64(bgd.fs.sb.InodePer_group) + subInodeNum,
	}
	inode.UpdateCsumAndWriteback()

	return inode
}

func(bgd *GroupDescriptor) GetFreeBlock() int64 {
	// Find free block in bitmap
	start := bgd.GetBlockBitmapLoc() * bgd.fs.sb.GetBlockSize()
	bgd.fs.dev.Seek(start, 0)
	subBlockNum := int64(-1)
	for i := 0; i < int(bgd.fs.sb.BlockPer_group/8); i++ {
		b := make([]byte, 1)
		bgd.fs.dev.Read(b)
		if b[0] != 0xFF {
			bitNum := bits.TrailingZeros8(^b[0])
			subBlockNum = int64(i) * 8 + int64(bitNum)
			b[0] |= 1 << uint(bitNum)
			bgd.fs.dev.Seek(-1, 1)
			bgd.fs.dev.Write(b)
			break
		}
	}

	if subBlockNum < 0 {
		return 0
	}

	// Update block bitmap checksum
	checksummer := NewChecksummer(bgd.fs.sb)
	checksummer.Write(bgd.fs.sb.Uuid[:])
	bgd.fs.dev.Seek(start, 0)
	b := make([]byte, int64(bgd.fs.sb.ClusterPer_group) / 8)
	bgd.fs.dev.Read(b)
	checksummer.Write(b)
	bgd.Block_bitmap_csum_lo = uint16(checksummer.Get() & 0xFFFF)
	bgd.Block_bitmap_csum_hi = uint16(checksummer.Get() >> 16)

	bgd.Free_blocks_count_lo--
	bgd.UpdateCsumAndWriteback()

	bgd.fs.sb.Free_blockCount_lo--
	bgd.fs.sb.UpdateCsumAndWriteback()

	return bgd.address / bgd.fs.sb.GetBlockSize() + subBlockNum - 1
}