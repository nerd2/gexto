package gexto

import (
	"github.com/lunixbochs/struc"
	"encoding/binary"
	"log"
	"bytes"
)

type Inode struct {
	Mode         int16     `struc:"int16,little"`
	Uid          int16     `struc:"int16,little"`
	Size_lo      int32     `struc:"int32,little"`
	Atime        int32     `struc:"int32,little"`
	Ctime        int32     `struc:"int32,little"`
	Mtime        int32     `struc:"int32,little"`
	Dtime        int32     `struc:"int32,little"`
	Gid          int16     `struc:"int16,little"`
	Links_count  int16     `struc:"int16,little"`
	Blocks_lo    int32     `struc:"int32,little"`
	Flags        int32     `struc:"int32,little"`
	Osd1         int32     `struc:"int32,little"`
	BlockOrExtents []byte `struc:"[60]byte,little"`
	Generation   int32     `struc:"int32,little"`
	File_acl_lo  int32     `struc:"int32,little"`
	Size_high    int32     `struc:"int32,little"`
	Obso_faddr   int32     `struc:"int32,little"`
	Osd2         [12]byte  `struc:"[12]byte"`
	Extra_isize  int16     `struc:"int16,little"`
	Checksum_hi  int16     `struc:"int16,little"`
	Ctime_extra  int32     `struc:"int32,little"`
	Mtime_extra  int32     `struc:"int32,little"`
	Atime_extra  int32     `struc:"int32,little"`
	Crtime       int32     `struc:"int32,little"`
	Crtime_extra int32     `struc:"int32,little"`
	Version_hi   int32     `struc:"int32,little"`
	Projid       int32     `struc:"int32,little"`
	fs           *fs
};


func (inode *Inode) UsesExtents() bool {
	return (inode.Flags & EXTENTS_FL) != 0
}

func (inode *Inode) UsesDirectoryHashTree() bool {
	return (inode.Flags & INDEX_FL) != 0
}

func (inode *Inode) ReadDirectory() []DirectoryEntry2 {
	sb := inode.fs.sb
	dev := inode.fs.dev

	if inode.UsesDirectoryHashTree() {
		log.Fatalf("Not implemented")
	}
	if inode.UsesExtents() {
		log.Fatalf("Not implemented")

		extentHeader := &ExtentHeader{}
		struc.Unpack(bytes.NewReader([]byte(inode.BlockOrExtents)), &extentHeader)
		log.Printf("extent header: %+v", extentHeader)
		if extentHeader.Depth == 0 { // Leaf
			for i := int16(0); i < extentHeader.Entries; i++ {
				extent := &Extent{}
				struc.Unpack(bytes.NewReader([]byte(inode.BlockOrExtents)[12 + i * 12:]), &extent)
				log.Printf("extent: %+v", extent)
			}
		} else {
			log.Fatalf("Not implemented")
		}
		return nil
	} else {
		ret := []DirectoryEntry2{}

		for blockTableIndex := int64(0); blockTableIndex < (inode.GetSize() + sb.GetBlockSize() - 1) / sb.GetBlockSize(); blockTableIndex++ {
			blockNum := inode.GetBlockPtr(blockTableIndex)
			blockStart := int64(blockNum) * sb.GetBlockSize()
			pos := blockStart
			for i := 0; i < 16; i++ {
				dev.Seek(pos, 0)
				dirEntry := DirectoryEntry2{}
				struc.Unpack(dev, &dirEntry)
				log.Printf("dirEntry %s: %+v", string(dirEntry.Name), dirEntry)
				pos += int64(dirEntry.Rec_len)
				ret = append(ret, dirEntry)
				if pos == blockStart+sb.GetBlockSize() {
					log.Printf("Reached end of block, next block")
					break
				} else if pos > blockStart + sb.GetBlockSize() {
					log.Fatalf("Unexpected overflow out of block when directory listing")
				}
			}
		}
		return ret
	}
}

func (inode *Inode) GetBlockPtr(num int64) int64 {
	if num < 12 {
		return int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*num:]))
	}

	num -= 12

	indirectsPerBlock := inode.fs.sb.GetBlockSize() / 4
	if num < indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*12:]))
		return inode.getIndirectBlockPtr(ptr, num)
	}
	num -= indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*13:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / indirectsPerBlock)
		return inode.getIndirectBlockPtr(l1, num % indirectsPerBlock)
	}

	num -= indirectsPerBlock * indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock * indirectsPerBlock {
		log.Println("Triple indirection")

		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*14:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / (indirectsPerBlock * indirectsPerBlock))
		l2 := inode.getIndirectBlockPtr(l1, (num / indirectsPerBlock) % indirectsPerBlock)
		return inode.getIndirectBlockPtr(l2, num % (indirectsPerBlock * indirectsPerBlock))
	}

	log.Fatalf("Exceeded maximum possible block count")
	return 0
}

func (inode *Inode) getIndirectBlockPtr(blockNum int64, offset int64) int64 {
	inode.fs.dev.Seek(blockNum * inode.fs.sb.GetBlockSize() + offset * 4, 0)
	x := make([]byte, 4)
	inode.fs.dev.Read(x)
	return int64(binary.LittleEndian.Uint32(x))
}

func (inode *Inode) GetSize() int64 {
	return (int64(inode.Size_high) << 32) | int64(inode.Size_lo)
}