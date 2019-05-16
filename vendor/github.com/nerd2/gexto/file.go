package gexto

import (
	"io"
	"fmt"
	"log"
)

type extFile struct {
	fs *fs
	inode *Inode
	pos int64
}

func (f *File) Read(p []byte) (n int, err error) {
	//log.Println("read", len(p), f.pos, f.inode.GetSize())
	blockNum := f.pos / f.fs.sb.GetBlockSize()
	blockPos := f.pos % f.fs.sb.GetBlockSize()
	len := int64(len(p))
	offset := int64(0)

	//log.Println("Read", len, f.pos, f.inode.GetSize())
	if len + f.pos > int64(f.inode.GetSize()) {
		len = int64(f.inode.GetSize()) - f.pos
	}

	if len <= 0 {
		//log.Println("EOF")
		return 0, io.EOF
	}

	for len > 0 {
		blockPtr, contiguousBlocks, found := f.inode.GetBlockPtr(blockNum)
		if !found {
			return int(offset), io.ErrUnexpectedEOF
		}

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

func (f *File) Write(p []byte) (n int, err error) {
	totalLen := len(p)

	//log.Println("Doing write", totalLen, p)

	for len(p) > 0 {
		blockNum := f.pos / f.fs.sb.GetBlockSize()
		blockPos := f.pos % f.fs.sb.GetBlockSize()

		//log.Println("Doing write", f.pos, blockNum, blockPos)

		blockPtr, contiguousBlocks, found := f.inode.GetBlockPtr(blockNum)

		if !found {
			//log.Println("Not found, extending")
			blockPtr, contiguousBlocks = f.inode.AddBlocks((int64(len(p)) + f.inode.fs.sb.GetBlockSize() - 1) / f.inode.fs.sb.GetBlockSize())
		}

		//log.Println(blockNum, blockPos, blockPtr, contiguousBlocks, len(p))
		writable := contiguousBlocks * f.fs.sb.GetBlockSize() - blockPos

		if writable == 0 {
			log.Fatalf("panic")
		}

		if writable > int64(len(p)) {
			writable = int64(len(p))
		}

		f.pos += writable
		//log.Println("seek", blockPtr * f.fs.sb.GetBlockSize() + blockPos, "write", writable)
		f.fs.dev.Seek(blockPtr * f.fs.sb.GetBlockSize() + blockPos, 0)
		f.fs.dev.Write(p[:writable])
		p = p[writable:]
	}

	if f.inode.GetSize() < f.pos {
		f.inode.SetSize(f.inode.GetSize() + int64(totalLen))
	}
	//log.Println("Write complete")

	return totalLen, nil
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
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
