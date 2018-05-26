package gexto

import (
	"io"
	"fmt"
)

type extFile struct {
	fs *fs
	inode *Inode
	pos int64
}

func (f *extFile) Read(p []byte) (n int, err error) {
	//log.Println("read", len(p), f.pos, f.inode.GetSize())
	blockNum := f.pos / f.fs.sb.GetBlockSize()
	blockPos := f.pos % f.fs.sb.GetBlockSize()
	len := int64(len(p))
	offset := int64(0)

	if len + f.pos > int64(f.inode.GetSize()) {
		len = int64(f.inode.GetSize()) - f.pos
	}

	if len == 0 {
		return 0, io.EOF
	}

	for len > 0 {
		blockPtr, contiguousBlocks := f.inode.GetBlockPtr(blockNum)

		f.fs.dev.Seek(blockPtr * f.fs.sb.GetBlockSize() + blockPos, 0)

		blockReadLen := contiguousBlocks * f.fs.sb.GetBlockSize() - blockPos
		if blockReadLen > len {
			blockReadLen = len
		}
		//log.Println(len, blockNum, blockPos, blockPtr, blockReadLen, offset)
		n, err := io.LimitReader(f.fs.dev, blockReadLen).Read(p[offset:])
		if err != nil {
			return 0, err
		}
		offset += int64(n)
		blockPos = 0
		blockNum++
		len -= int64(n)
	}
	f.pos += offset
	//log.Println(int(offset))
	return int(offset), nil
}

func (f *extFile) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
		f.pos = offset
	case 1:
		f.pos += offset
	case 2:
		f.pos = f.inode.GetSize() - offset
	default:
		return 0, fmt.Errorf("Unsupported whence")
	}

	if f.pos >= f.inode.GetSize() {
		return f.inode.GetSize(), io.EOF
	} else {
		return f.pos, nil
	}
}