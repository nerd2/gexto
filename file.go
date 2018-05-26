package gexto

import (
	"io"
	"encoding/binary"
	"log"
)

type extFile struct {
	fs *fs
	inode *Inode
	pos int64
}

func (f *extFile) Read(p []byte) (n int, err error) {
	log.Println("read", len(p))
	blockNum := f.pos / f.fs.sb.GetBlockSize()
	if blockNum > 11 {
		log.Fatalf("Not implemented")
	}
	blockPos := f.pos % f.fs.sb.GetBlockSize()
	len := int64(len(p))
	offset := int64(0)

	var endErr error
	if len + f.pos > int64(f.inode.Size_lo) {
		len = int64(f.inode.Size_lo) - f.pos
		endErr = io.EOF
	}

	for len > 0 {
		blockPtr := int64(binary.LittleEndian.Uint32(f.inode.BlockOrExtents[4*blockNum:]))

		f.fs.dev.Seek(blockPtr * f.fs.sb.GetBlockSize() + blockPos, 0)

		blockReadLen := f.fs.sb.GetBlockSize() - blockPos
		if blockReadLen > len {
			blockReadLen = len
		}
		n, err := io.LimitReader(f.fs.dev, blockReadLen).Read(p[offset:])
		if err != nil {
			return 0, err
		}
		offset += int64(n)
		blockPos = 0
		blockNum++
		len -= int64(n)
	}
	log.Println(int(offset))
	return int(offset), endErr
}