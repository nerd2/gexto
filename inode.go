package gexto

import (
	"github.com/lunixbochs/struc"
	"os"
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

func (inode *Inode) ReadFile(sb *Superblock, dev *os.File) {
	size := int64(inode.Size_lo)
	for blockTableIndex := int64(0); blockTableIndex < (int64(inode.Size_lo)+sb.GetBlockSize()-1)/sb.GetBlockSize(); blockTableIndex++ {
		blockNum := binary.LittleEndian.Uint32(inode.BlockOrExtents[blockTableIndex * 4:])
		dev.Seek(int64(blockNum) * sb.GetBlockSize(), 0)
		sizeInBlock := sb.GetBlockSize()
		if size < sizeInBlock {
			sizeInBlock = size
		}
		data := make([]byte, sizeInBlock)
		dev.Read(data)
		log.Printf("%s", string(data))
		size -= sizeInBlock
	}

	if size > 0 {
		log.Fatalf("Oversize block")
	}
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

		for blockTableIndex := int64(0); blockTableIndex < (int64(inode.Size_lo) + sb.GetBlockSize() - 1) / sb.GetBlockSize(); blockTableIndex++ {
			blockNum := binary.LittleEndian.Uint32(inode.BlockOrExtents[blockTableIndex * 4:])
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

