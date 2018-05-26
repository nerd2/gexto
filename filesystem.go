package gexto

import (
	"os"
	"strings"
	"github.com/lunixbochs/struc"
	"log"
	"fmt"
	"io"
)

type fs struct {
	sb *Superblock
	dev *os.File
}

func (fs *fs) Open(name string) (*File, error) {
	parts := strings.Split(name, "/")

	inodeNum := int64(ROOT_INO)
	var inode *Inode
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		log.Println(part)
		inode = fs.getInode(inodeNum)
		dirContents := inode.ReadDirectory()
		found := false
		for i := 0; i < len(dirContents); i++ {
			log.Println(string(dirContents[i].Name), part, dirContents[i].Flags, dirContents[i].Inode)
			if string(dirContents[i].Name) == part {
				found = true
				inodeNum = int64(dirContents[i].Inode)
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("No such file or directory")
		}
	}

	return &File{extFile{
		fs: fs,
		inode: fs.getInode(inodeNum),
		pos: 0,
	}}, nil
}

func (fs *fs) Create(name string) (*File, error) {
	return nil, nil
}

func (fs *fs) Remove(name string) error {
	return nil
}

func (fs *fs) Mkdir(name string, perm os.FileMode) error {
	return nil
}

// --------------------------


func (fs *fs) getInode(inodeAddress int64) *Inode {
	bgd := fs.getBlockGroupDescriptor((inodeAddress - 1) / int64(fs.sb.InodePer_group))
	index := (inodeAddress - 1) % int64(fs.sb.InodePer_group)
	pos := bgd.GetInodeTableLoc(fs.sb) * fs.sb.GetBlockSize() + index * int64(fs.sb.Inode_size)
	log.Printf("%d %d %d %d", bgd.GetInodeTableLoc(fs.sb), fs.sb.GetBlockSize(), index, fs.sb.Inode_size)
	fs.dev.Seek(pos, 0)

	inode := &Inode{fs:fs}
	struc.Unpack(fs.dev, &inode)
	log.Printf("Read inode at offset %d, contents:\n%+v\n", pos, inode)
	return inode
}

func (fs *fs) getBlockGroupDescriptor(blockGroupNum int64) *GroupDescriptor {
	blockSize := fs.sb.GetBlockSize()
	bgdtLocation := 1024/blockSize + 1

	log.Println(bgdtLocation,blockGroupNum,blockSize)
	bgd := &GroupDescriptor{}
	if fs.sb.FeatureIncompat64bit() {
		fs.dev.Seek(bgdtLocation*blockSize + 64 * blockGroupNum, 0)
		struc.Unpack(fs.dev, &bgd)
	} else {
		fs.dev.Seek(bgdtLocation*blockSize + 32 * blockGroupNum, 0)
		struc.Unpack(io.LimitReader(fs.dev, 32), &bgd)
	}
	log.Printf("Read block group %d, contents:\n%+v\n", blockGroupNum, bgd)
	return bgd
}