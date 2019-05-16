package gexto

import (
	"os"
	"github.com/lunixbochs/struc"
	"fmt"
	"syscall"
)

type File struct {
	extFile
}

type FileSystem interface {
	Open(name string) (*File, error)
	Create(name string) (*File, error)
	Remove(name string) error
	Mkdir(name string, perm os.FileMode) error
	Close() error
}

func NewFileSystem(devicePath string) (FileSystem, error) {
	f, err := os.OpenFile(devicePath, syscall.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	ret := fs{}

	f.Seek(1024, 0)

	ret.dev = f
	ret.sb = &Superblock{
		address: 1024,
		fs: &ret,
	}
	err = struc.Unpack(f, ret.sb)
	if err != nil {
		return nil, err
	}

	//log.Printf("Super:\n%+v\n", *ret.sb)

	numBlockGroups := (ret.sb.GetBlockCount() + int64(ret.sb.BlockPer_group) - 1) / int64(ret.sb.BlockPer_group)
	numBlockGroups2 := (ret.sb.InodeCount + ret.sb.InodePer_group - 1) / ret.sb.InodePer_group
	if numBlockGroups != int64(numBlockGroups2) {
		return nil, fmt.Errorf("Block/inode mismatch: %d %d %d", ret.sb.GetBlockCount(), numBlockGroups, numBlockGroups2)
	}

	ret.sb.numBlockGroups = numBlockGroups

	return &ret, nil
}


