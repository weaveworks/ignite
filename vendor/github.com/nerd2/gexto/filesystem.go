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

		inode = fs.getInode(inodeNum)
		dirContents := inode.ReadDirectory()
		found := false
		for i := 0; i < len(dirContents); i++ {
			//log.Println(string(dirContents[i].Name), part, dirContents[i].Flags, dirContents[i].Inode)
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

	inode = fs.getInode(inodeNum)
	//log.Printf("Inode %d with mode %x", inode.num, inode.Mode)
	return &File{extFile{
		fs: fs,
		inode: inode,
		pos: 0,
	}}, nil
}

func (fs *fs) Create(path string) (*File, error) {
	log.Println("CREATE", path)
	parts := strings.Split(path, "/")

	inode := fs.getInode(int64(ROOT_INO))

	for _, part := range parts[:len(parts)-1] {
		if len(part) == 0 {
			continue
		}

		dirContents := inode.ReadDirectory()
		found := false
		for i := 0; i < len(dirContents); i++ {
			//log.Println(string(dirContents[i].Name), part, dirContents[i].Flags, dirContents[i].Inode)
			if string(dirContents[i].Name) == part {
				found = true
				inode = fs.getInode(int64(dirContents[i].Inode))
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("No such file or directory")
		}
	}

	name := parts[len(parts)-1]

	newFile := fs.CreateNewFile(0777)
	log.Printf("Creating new file with inode %d and perms %d", newFile.inode.num, newFile.inode.Mode)
	newFile.inode.Mode |= 0x8000
	newFile.inode.UpdateCsumAndWriteback()

	NewDirectory(inode).AddEntry(&DirectoryEntry2{
		Inode:    uint32(newFile.inode.num),
		Flags:    0,
		Name:     name,
	})

	return newFile, nil
}

func (fs *fs) Remove(name string) error {
	return nil
}

func (fs *fs) Mkdir(path string, perm os.FileMode) error {
	log.Println("MKDIR", path)
	parts := strings.Split(path, "/")

	inode := fs.getInode(int64(ROOT_INO))

	for _, part := range parts[:len(parts)-1] {
		if len(part) == 0 {
			continue
		}

		dirContents := inode.ReadDirectory()
		found := false
		for i := 0; i < len(dirContents); i++ {
			//log.Println(string(dirContents[i].Name), part, dirContents[i].Flags, dirContents[i].Inode)
			if string(dirContents[i].Name) == part {
				found = true
				inode = fs.getInode(int64(dirContents[i].Inode))
				break
			}
		}

		if !found {
			return fmt.Errorf("No such file or directory")
		}
	}

	name := parts[len(parts)-1]

	newFile := fs.CreateNewFile(perm)
	log.Printf("Creating new directory with inode %d and perms %d", newFile.inode.num, newFile.inode.Mode)
	newFile.inode.Mode |= 0x4000
	newFile.inode.UpdateCsumAndWriteback()

	{
		checksummer := NewChecksummer(inode.fs.sb)
		checksummer.Write(inode.fs.sb.Uuid[:])
		checksummer.WriteUint32(uint32(newFile.inode.num))
		checksummer.WriteUint32(uint32(newFile.inode.Generation))

		dirEntryDot := DirectoryEntry2{
			Inode:    uint32(newFile.inode.num),
			Flags:    2,
			Rec_len:  12,
			Name:     ".",
		}
		recLenDot, _ := struc.Sizeof(&dirEntryDot)
		struc.Pack(checksummer, dirEntryDot)
		struc.Pack(newFile, dirEntryDot)
		{
			blank1 := make([]byte, 12-recLenDot)
			checksummer.Write(blank1)
			newFile.Write(blank1)
		}

		dirEntryDotDot := DirectoryEntry2{
			Inode:    uint32(inode.num),
			Flags:    2,
			Name:     "..",
		}
		recLenDotDot, _ := struc.Sizeof(&dirEntryDotDot)
		dirEntryDotDot.Rec_len = uint16(1024 - 12 - 12)
		struc.Pack(checksummer, dirEntryDotDot)
		struc.Pack(newFile, dirEntryDotDot)

		blank := make([]byte, 1024 - 12 - 12 - recLenDotDot)
		checksummer.Write(blank)
		newFile.Write(blank)

		dirSum := DirectoryEntryCsum{
			FakeInodeZero: 0,
			Rec_len:  uint16(12),
			FakeName_len: 0,
			FakeFileType:    0xDE,
			Checksum:     checksummer.Get(),
		}
		struc.Pack(newFile, &dirSum)
	}

	NewDirectory(inode).AddEntry(&DirectoryEntry2{
		Inode:    uint32(newFile.inode.num),
		Flags:    0,
		Name:     name,
	})

	newFile.inode.Links_count++
	newFile.inode.UpdateCsumAndWriteback()

	inode.Links_count++
	inode.UpdateCsumAndWriteback()

	bgd := fs.getBlockGroupDescriptor((newFile.inode.num-1) / int64(inode.fs.sb.InodePer_group))
	bgd.Used_dirs_count_lo++
	bgd.UpdateCsumAndWriteback()

	return nil
}

func (fs *fs) Close() error {
	err := fs.dev.Close()
	if err != nil {
		return err
	}
	fs.sb = nil
	fs.dev = nil
	return nil
}

// --------------------------


func (fs *fs) getInode(inodeAddress int64) *Inode {
	bgd := fs.getBlockGroupDescriptor((inodeAddress - 1) / int64(fs.sb.InodePer_group))
	index := (inodeAddress - 1) % int64(fs.sb.InodePer_group)
	pos := bgd.GetInodeTableLoc() * fs.sb.GetBlockSize() + index * int64(fs.sb.Inode_size)
	//log.Printf("%d %d %d %d", bgd.GetInodeTableLoc(), fs.sb.GetBlockSize(), index, fs.sb.Inode_size)
	fs.dev.Seek(pos, 0)

	inode := &Inode{
		fs: fs,
		address: pos,
		num: inodeAddress,}
	struc.Unpack(fs.dev, &inode)
	//log.Printf("Read inode %d, contents:\n%+v\n", inodeAddress, inode)
	return inode
}

func (fs *fs) getBlockGroupDescriptor(blockGroupNum int64) *GroupDescriptor {
	blockSize := fs.sb.GetBlockSize()
	bgdtLocation := 1024/blockSize + 1

	size := int64(32)
	if fs.sb.FeatureIncompat64bit() {
		size = int64(64)
	}
	addr := bgdtLocation*blockSize + size * blockGroupNum
	bgd := &GroupDescriptor{
		fs:fs,
		address: addr,
		num: blockGroupNum,
	}
	fs.dev.Seek(addr, 0)
	struc.Unpack(io.LimitReader(fs.dev, size), &bgd)
	//log.Printf("Read block group %d, contents:\n%+v\n", blockGroupNum, bgd)
	return bgd
}

func (fs *fs) CreateNewFile(perm os.FileMode) *File {
	var inode *Inode
	for i := int64(0); i < fs.sb.numBlockGroups; i++ {
		bgd := fs.getBlockGroupDescriptor(i)
		inode = bgd.GetFreeInode()
		if inode != nil {
			break
		}
	}

	if inode == nil {
		log.Fatalln("Couldn't get free inode", fs.sb.numBlockGroups, fs.sb.Free_inodeCount)
		return nil
	}

	inode.Mode = uint16(perm & 0x1FF)
	inode.UpdateCsumAndWriteback()

	return &File{extFile{
		fs: fs,
		inode: inode,
	}}
}

func (fs *fs) GetFreeBlocks(n int) (int64, int64) {
	for i := int64(0); i < fs.sb.numBlockGroups; i++ {
		bgd := fs.getBlockGroupDescriptor(i)
		blockNum, numBlocks := bgd.GetFreeBlocks(int64(n))
		if blockNum > 0 {
			return blockNum + i * int64(fs.sb.BlockPer_group), numBlocks
		}
	}
	log.Fatalf("Failed to find free block")
	return 0, 0
}