package gexto

import (
	"github.com/lunixbochs/struc"
	"encoding/binary"
	"log"
	"io"
	"bytes"
)

type MoveExtent struct {
	Reserved    uint32  `struc:"uint32,little"`
	Donor_fd    uint32  `struc:"uint32,little"`
	Orig_start  uint64 `struc:"uint64,little"`
	Donor_start uint64 `struc:"uint64,little"`
	Len         uint64 `struc:"uint64,little"`
	Moved_len   uint64 `struc:"uint64,little"`
};

type ExtentHeader struct {
	Magic      uint16 `struc:"uint16,little"`
	Entries    uint16 `struc:"uint16,little"`
	Max        uint16 `struc:"uint16,little"`
	Depth      uint16 `struc:"uint16,little"`
	Generation uint32 `struc:"uint32,little"`
}

type ExtentInternal struct {
	Block     uint32 `struc:"uint32,little"`
	Leaf_low  uint32 `struc:"uint32,little"`
	Leaf_high uint16 `struc:"uint16,little"`
	Unused    uint16 `struc:"uint16,little"`
}

type Extent struct {
	Block    uint32 `struc:"uint32,little"`
	Len      uint16 `struc:"uint16,little"`
	Start_hi uint16 `struc:"uint16,little"`
	Start_lo uint32 `struc:"uint32,little"`
}

type DirectoryEntry2 struct {
	Inode uint32 `struc:"uint32,little"`
	Rec_len uint16 `struc:"uint16,little"`
	Name_len uint8 `struc:"uint8,sizeof=Name"`
	Flags uint8 `struc:"uint8"`
	Name string `struc:"[]byte"`
}

type Inode struct {
	Mode           uint16   `struc:"uint16,little"`
	Uid            uint16   `struc:"uint16,little"`
	Size_lo        uint32   `struc:"uint32,little"`
	Atime          uint32   `struc:"uint32,little"`
	Ctime          uint32   `struc:"uint32,little"`
	Mtime          uint32   `struc:"uint32,little"`
	Dtime          uint32   `struc:"uint32,little"`
	Gid            uint16   `struc:"uint16,little"`
	Links_count    uint16   `struc:"uint16,little"`
	Blocks_lo      uint32   `struc:"uint32,little"`
	Flags          uint32   `struc:"uint32,little"`
	Osd1           uint32   `struc:"uint32,little"`
	BlockOrExtents []byte   `struc:"[60]byte,little"`
	Generation     uint32   `struc:"uint32,little"`
	File_acl_lo    uint32   `struc:"uint32,little"`
	Size_high      uint32   `struc:"uint32,little"`
	Obso_faddr     uint32   `struc:"uint32,little"`
	Osd2           [12]byte `struc:"[12]byte"`
	Extra_isize    uint16   `struc:"uint16,little"`
	Checksum_hi    uint16   `struc:"uint16,little"`
	Ctime_extra    uint32   `struc:"uint32,little"`
	Mtime_extra    uint32   `struc:"uint32,little"`
	Atime_extra    uint32   `struc:"uint32,little"`
	Crtime         uint32   `struc:"uint32,little"`
	Crtime_extra   uint32   `struc:"uint32,little"`
	Version_hi     uint32   `struc:"uint32,little"`
	Projid         uint32   `struc:"uint32,little"`
	fs             *fs
};


func (inode *Inode) UsesExtents() bool {
	return (inode.Flags & EXTENTS_FL) != 0
}

func (inode *Inode) UsesDirectoryHashTree() bool {
	return (inode.Flags & INDEX_FL) != 0
}

func (inode *Inode) ReadDirectory() []DirectoryEntry2 {
	if inode.UsesDirectoryHashTree() {
		log.Fatalf("Not implemented")
	}

	f := &File{extFile{
		fs: inode.fs,
		inode: inode,
		pos: 0,
	}}

	ret := []DirectoryEntry2{}
	for {
		start, _ := f.Seek(0, 1)
		dirEntry := DirectoryEntry2{}
		err := struc.Unpack(f, &dirEntry)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf(err.Error())
		}
		//log.Printf("dirEntry %s: %+v", string(dirEntry.Name), dirEntry)
		f.Seek(int64(dirEntry.Rec_len) + start, 0)
		ret = append(ret, dirEntry)
	}
	return ret
}

// Returns the blockId of the file block, and the number of contiguous blocks
func (inode *Inode) GetBlockPtr(num int64) (int64, int64) {
	if inode.UsesExtents() {
		log.Println("Finding", num)
		r := io.Reader(bytes.NewReader(inode.BlockOrExtents))

		for {
			extentHeader := &ExtentHeader{}
			struc.Unpack(r, &extentHeader)
			log.Printf("extent header: %+v", extentHeader)
			if extentHeader.Depth == 0 { // Leaf
				for i := uint16(0); i < extentHeader.Entries; i++ {
					extent := &Extent{}
					struc.Unpack(r, &extent)
					log.Printf("extent leaf: %+v", extent)
					if int64(extent.Block) <= num && int64(extent.Block)+int64(extent.Len) > num {
						log.Println("Found")
						return int64(extent.Start_hi<<32) + int64(extent.Start_lo) + num - int64(extent.Block), int64(extent.Block) + int64(extent.Len) - num
					}
				}
				log.Fatalf("Could not find in extent tree")
				return 0,0
			} else {
				found := false
				for i := uint16(0); i < extentHeader.Entries; i++ {
					extent := &ExtentInternal{}
					struc.Unpack(r, &extent)
					log.Printf("extent internal: %+v", extent)
					if int64(extent.Block) <= num {
						newBlock := int64(extent.Leaf_high<<32) + int64(extent.Leaf_low)
						inode.fs.dev.Seek(newBlock * inode.fs.sb.GetBlockSize(), 0)
						r = inode.fs.dev
						found = true
						break
					}
				}
				if !found {
					log.Fatalf("Could not find in extent tree")
					return 0,0
				}
			}
		}

	}

	if num < 12 {
		return int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*num:])),1
	}

	num -= 12

	indirectsPerBlock := inode.fs.sb.GetBlockSize() / 4
	if num < indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*12:]))
		return inode.getIndirectBlockPtr(ptr, num),1
	}
	num -= indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*13:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / indirectsPerBlock)
		return inode.getIndirectBlockPtr(l1, num % indirectsPerBlock),1
	}

	num -= indirectsPerBlock * indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock * indirectsPerBlock {
		log.Println("Triple indirection")

		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*14:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / (indirectsPerBlock * indirectsPerBlock))
		l2 := inode.getIndirectBlockPtr(l1, (num / indirectsPerBlock) % indirectsPerBlock)
		return inode.getIndirectBlockPtr(l2, num % (indirectsPerBlock * indirectsPerBlock)),1
	}

	log.Fatalf("Exceeded maximum possible block count")
	return 0,0
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