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

type DirectoryEntryCsum struct {
	FakeInodeZero uint32 `struc:"uint32,little"`
	Rec_len uint16 `struc:"uint16,little"`
	FakeName_len uint8 `struc:"uint8"`
	FakeFileType uint8 `struc:"uint8"`
	Checksum uint32 `struc:"uint32,little"`
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
	BlockOrExtents [60]byte `struc:"[60]byte,little"`
	Generation     uint32   `struc:"uint32,little"`
	File_acl_lo    uint32   `struc:"uint32,little"`
	Size_high      uint32   `struc:"uint32,little"`
	Obso_faddr     uint32   `struc:"uint32,little"`
	// OSD2 - linux only starts
	Blocks_high    uint16   `struc:"uint16,little"`
	File_acl_high  uint16   `struc:"uint16,little"`
	Uid_high       uint16   `struc:"uint16,little"`
	Gid_high       uint16   `struc:"uint16,little"`
	Checksum_low   uint16   `struc:"uint16,little"`
	Unused         uint16   `struc:"uint16,little"`
	// OSD2 - linux only ends
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
	address        int64
	num            int64
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
		if dirEntry.Rec_len < 9 {
			log.Fatalf("corrupt direntry")
		}
		ret = append(ret, dirEntry)
	}
	return ret
}

func (inode *Inode) AddBlocks(n int64) (blockNum int64, contiguousBlocks int64) {
	if !inode.UsesExtents() {
		log.Fatalf("Not implemented")
	}

	r := inode.fs.dev
	r.Seek(inode.address + 40, 0)

	for {
		headerPos, _ := r.Seek(0,1)
		extentHeader := &ExtentHeader{}
		struc.Unpack(r, &extentHeader)
		//log.Printf("extent header: %+v", extentHeader)
		if extentHeader.Depth == 0 { // Leaf
			max := int64(0)
			for i := uint16(0); i < extentHeader.Entries; i++ {
				extent := &Extent{}
				struc.Unpack(r, &extent)
				upper := int64(extent.Block) + int64(extent.Len)
				if upper > max {
					max = upper
				}
			}
			if extentHeader.Entries < extentHeader.Max {
				savePos, _ := r.Seek(0, 1)
				blockNum, numBlocks := inode.fs.GetFreeBlocks(int(n))
				newExtent := &Extent{
					Block: uint32(max),
					Len: uint16(numBlocks),
					Start_hi: uint16(blockNum >> 32),
					Start_lo: uint32(blockNum & 0xFFFFFFFF),
				}
				r.Seek(savePos, 0)
				struc.Pack(r, &newExtent)
				extentHeader.Entries++
				//log.Println("Extended to", extentHeader.Entries, headerPos)
				r.Seek(headerPos, 0)
				struc.Pack(r, extentHeader)
				r.Seek(inode.address, 0)
				struc.Unpack(r, inode)
				inode.Blocks_lo += uint32(numBlocks*inode.fs.sb.GetBlockSize()/512)
				inode.UpdateCsumAndWriteback()

				//log.Println("AddBlocks", n, numBlocks)

				return blockNum, numBlocks
			} else {
				log.Fatalf("Unable to extend no room")
			}
		} else {
			max := uint32(0)
			var best *ExtentInternal
			for i := uint16(0); i < extentHeader.Entries; i++ {
				extent := &ExtentInternal{}
				struc.Unpack(r, &extent)
				//log.Printf("extent internal: %+v", extent)
				if extent.Block > max {
					best = extent
				}
			}

			newBlock := int64(best.Leaf_high<<32) + int64(best.Leaf_low)
			r.Seek(newBlock*inode.fs.sb.GetBlockSize(), 0)
		}
	}

	//log.Println("AddBlocks", n, 0)
	return 0,0
}

func (inode *Inode) UpdateCsumAndWriteback() {
	if inode.fs.sb.Inode_size != 128 {
		log.Fatalln("Unsupported inode size", inode.fs.sb.Inode_size)
	}

	cs := NewChecksummer(inode.fs.sb)

	cs.Write(inode.fs.sb.Uuid[:])
	cs.WriteUint32(uint32(inode.num))
	cs.WriteUint32(uint32(inode.Generation))
	inode.Checksum_low = 0
	struc.Pack(LimitWriter(cs, 128), inode)
	inode.Checksum_low = uint16(cs.Get() & 0xFFFF)

	inode.fs.dev.Seek(inode.address, 0)
	struc.Pack(LimitWriter(inode.fs.dev, 128), inode)
}

// Returns the blockId of the file block, and the number of contiguous blocks
func (inode *Inode) GetBlockPtr(num int64) (int64, int64, bool) {
	if inode.UsesExtents() {
		//log.Println("Finding", num)
		r := io.Reader(bytes.NewReader(inode.BlockOrExtents[:]))

		for {
			extentHeader := &ExtentHeader{}
			struc.Unpack(r, &extentHeader)
			//log.Printf("extent header: %+v", extentHeader)
			if extentHeader.Depth == 0 { // Leaf
				for i := uint16(0); i < extentHeader.Entries; i++ {
					extent := &Extent{}
					struc.Unpack(r, &extent)
					//log.Printf("extent leaf: %+v", extent)
					if int64(extent.Block) <= num && int64(extent.Block)+int64(extent.Len) > num {
						//log.Println("Found")
						return int64(extent.Start_hi<<32) + int64(extent.Start_lo) + num - int64(extent.Block), int64(extent.Block) + int64(extent.Len) - num, true
					}
				}
				return 0, 0, false
			} else {
				found := false
				for i := uint16(0); i < extentHeader.Entries; i++ {
					extent := &ExtentInternal{}
					struc.Unpack(r, &extent)
					//log.Printf("extent internal: %+v", extent)
					if int64(extent.Block) <= num {
						newBlock := int64(extent.Leaf_high<<32) + int64(extent.Leaf_low)
						inode.fs.dev.Seek(newBlock * inode.fs.sb.GetBlockSize(), 0)
						r = inode.fs.dev
						found = true
						break
					}
				}
				if !found {
					return 0,0, false
				}
			}
		}

	}

	if num < 12 {
		return int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*num:])), 1, true
	}

	num -= 12

	indirectsPerBlock := inode.fs.sb.GetBlockSize() / 4
	if num < indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*12:]))
		return inode.getIndirectBlockPtr(ptr, num),1, true
	}
	num -= indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock {
		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*13:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / indirectsPerBlock)
		return inode.getIndirectBlockPtr(l1, num % indirectsPerBlock),1, true
	}

	num -= indirectsPerBlock * indirectsPerBlock

	if num < indirectsPerBlock * indirectsPerBlock * indirectsPerBlock {
		log.Println("Triple indirection")

		ptr := int64(binary.LittleEndian.Uint32(inode.BlockOrExtents[4*14:]))
		l1 := inode.getIndirectBlockPtr(ptr, num / (indirectsPerBlock * indirectsPerBlock))
		l2 := inode.getIndirectBlockPtr(l1, (num / indirectsPerBlock) % indirectsPerBlock)
		return inode.getIndirectBlockPtr(l2, num % (indirectsPerBlock * indirectsPerBlock)),1, true
	}

	log.Fatalf("Exceeded maximum possible block count")
	return 0,0,false
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

func (inode *Inode) SetSize(i int64) {
	inode.Size_high = uint32(i >> 32)
	inode.Size_lo = uint32(i & 0xFFFFFFFF)
	inode.UpdateCsumAndWriteback()
}