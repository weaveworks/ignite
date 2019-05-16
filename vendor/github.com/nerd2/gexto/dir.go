package gexto

import (
	"github.com/lunixbochs/struc"
	"fmt"
	"log"
)

type directory struct {
	sb *Superblock
	f  *File
}

func NewDirectory(inode *Inode) *directory {
	return &directory{
		f: &File{
			extFile{
				fs:    inode.fs,
				inode: inode,
				pos:   0,
			},
		},
		sb: inode.fs.sb,
	}
}

func (dir *directory) AddEntry(entry *DirectoryEntry2) error {
	entrySize, _ := struc.Sizeof(entry)
	entry.Rec_len = uint16(entrySize)

	pos, _ := dir.f.Seek(0, 2)
	if pos % dir.sb.GetBlockSize() != 0 {
		return fmt.Errorf("Unexpected size of directory file: %d", pos)
	} else if pos == 0 {
		return fmt.Errorf("Unexpected empty directory")
	}
	dir.f.Seek(pos - dir.sb.GetBlockSize(), 0)

	//log.Println("AddEntry", pos)

	checksummer := NewChecksummer(dir.sb)
	checksummer.Write(dir.f.inode.fs.sb.Uuid[:])
	checksummer.WriteUint32(uint32(dir.f.inode.num))
	checksummer.WriteUint32(uint32(dir.f.inode.Generation))

	totalLen := int64(0)
	modified := false
	for totalLen < dir.sb.GetBlockSize() {
		//log.Println("AddEntry loop ", totalLen)
		if totalLen == dir.sb.GetBlockSize() - 12 {
			//log.Println("AddEntry found checksum", modified)
			if modified {
				dirSum := DirectoryEntryCsum{
					FakeInodeZero: 0,
					Rec_len:  uint16(12),
					FakeName_len: 0,
					FakeFileType:    0xDE,
					Checksum:     checksummer.Get(),
				}
				struc.Pack(dir.f, &dirSum)
			}
			break
		}

		dirEntry := &DirectoryEntry2{}
		err := struc.Unpack(dir.f, dirEntry)
		if err != nil {
			return err
		}

		//log.Println("AddEntry found entry", dirEntry.Rec_len, dirEntry.Name)

		if dirEntry.Rec_len == 0 {
			log.Fatalf("Invalid")
		}

		deSize, _ := struc.Sizeof(dirEntry)
		if !modified && int64(dirEntry.Rec_len) >= int64(deSize) + int64(entrySize) {
			//log.Println("Found a hole", dirEntry.Rec_len, deSize, entrySize)
			dir.f.Seek(pos - dir.sb.GetBlockSize() + int64(totalLen), 0)
			newDeSize := (deSize + 3) & ^3
			entry.Rec_len = dirEntry.Rec_len - uint16(newDeSize)
			dirEntry.Rec_len = uint16(newDeSize)
			struc.Pack(dir.f, dirEntry)
			struc.Pack(checksummer, dirEntry)
			pad1 := make([]byte, newDeSize - deSize)
			dir.f.Write(pad1)
			checksummer.Write(pad1)
			struc.Pack(dir.f, entry)
			struc.Pack(checksummer, entry)
			pad := make([]byte, int(entry.Rec_len) - entrySize)
			dir.f.Write(pad)
			checksummer.Write(pad)
			totalLen += int64(dirEntry.Rec_len) + int64(entry.Rec_len)
			modified = true
		} else {
			struc.Pack(checksummer, dirEntry)
			skip := int64(dirEntry.Rec_len) - int64(deSize)
			dir.f.Seek(skip, 1)
			checksummer.Write(make([]byte, skip))
			totalLen += int64(dirEntry.Rec_len)
		}
	}

	if !modified {
		log.Println("No hole found, adding new block")
		checksummer := NewChecksummer(dir.sb)
		checksummer.Write(dir.f.inode.fs.sb.Uuid[:])
		checksummer.WriteUint32(uint32(dir.f.inode.num))
		checksummer.WriteUint32(uint32(dir.f.inode.Generation))

		entry.Rec_len = uint16(dir.sb.GetBlockSize() - 12)
		struc.Pack(dir.f, entry)
		struc.Pack(checksummer, entry)
		dirSum := DirectoryEntryCsum{
			FakeInodeZero: 0,
			Rec_len:  uint16(12),
			FakeName_len: 0,
			FakeFileType:    0xDE,
			Checksum:     checksummer.Get(),
		}
		struc.Pack(dir.f, &dirSum)
	}

	return nil
}
