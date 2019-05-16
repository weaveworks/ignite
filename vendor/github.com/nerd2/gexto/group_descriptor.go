package gexto

import (
	"github.com/lunixbochs/struc"
	"math/bits"
)

type GroupDescriptor struct {
	Block_bitmap_lo      uint32 `struc:"uint32,little"`
	Inode_bitmap_lo      uint32 `struc:"uint32,little"`
	Inode_table_lo       uint32 `struc:"uint32,little"`
	Free_blocks_count_lo uint16 `struc:"uint16,little"`
	Free_inodes_count_lo uint16 `struc:"uint16,little"`
	Used_dirs_count_lo   uint16 `struc:"uint16,little"`
	Flags                uint16 `struc:"uint16,little"`
	Exclude_bitmap_lo    uint32 `struc:"uint32,little"`
	Block_bitmap_csum_lo uint16 `struc:"uint16,little"`
	Inode_bitmap_csum_lo uint16 `struc:"uint16,little"`
	Itable_unused_lo     uint16 `struc:"uint16,little"`
	Checksum             uint16 `struc:"uint16,little"`
	Block_bitmap_hi      uint32 `struc:"uint32,little"`
	Inode_bitmap_hi      uint32 `struc:"uint32,little"`
	Inode_table_hi       uint32 `struc:"uint32,little"`
	Free_blocks_count_hi uint16 `struc:"uint16,little"`
	Free_inodes_count_hi uint16 `struc:"uint16,little"`
	Used_dirs_count_hi   uint16 `struc:"uint16,little"`
	Itable_unused_hi     uint16 `struc:"uint16,little"`
	Exclude_bitmap_hi    uint32 `struc:"uint32,little"`
	Block_bitmap_csum_hi uint16 `struc:"uint16,little"`
	Inode_bitmap_csum_hi uint16 `struc:"uint16,little"`
	Reserved             uint32 `struc:"uint32,little"`
	fs                   *fs
	num                  int64
	address              int64
};

func (bgd *GroupDescriptor) GetInodeBitmapLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_bitmap_hi) << 32) | int64(bgd.Inode_bitmap_lo)
	} else {
		return int64(bgd.Inode_bitmap_lo)
	}
}

func (bgd *GroupDescriptor) GetInodeTableLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Inode_table_hi) << 32) | int64(bgd.Inode_table_lo)
	} else {
		return int64(bgd.Inode_table_lo)
	}
}

func (bgd *GroupDescriptor) GetBlockBitmapLoc() int64 {
	if bgd.fs.sb.FeatureIncompat64bit() {
		return (int64(bgd.Block_bitmap_hi) << 32) | int64(bgd.Block_bitmap_lo)
	} else {
		return int64(bgd.Block_bitmap_lo)
	}
}

func (bgd *GroupDescriptor) UpdateCsumAndWriteback() {
	cs := NewChecksummer(bgd.fs.sb)

	cs.Write(bgd.fs.sb.Uuid[:])
	cs.WriteUint32(uint32(bgd.num))
	bgd.Checksum = 0
	struc.Pack(cs, bgd)
	bgd.Checksum = uint16(cs.Get() & 0xFFFF)

	bgd.fs.dev.Seek(bgd.address, 0)
	struc.Pack(bgd.fs.dev, bgd)
}

func(bgd *GroupDescriptor) GetFreeInode() *Inode {
	start := bgd.GetInodeBitmapLoc() * bgd.fs.sb.GetBlockSize()
	bgd.fs.dev.Seek(start, 0)

	subInodeNum := int64(-1)

	if bgd.Flags & BG_INODE_UNINIT != 0 {
		b := make([]byte, bgd.fs.sb.InodePer_group/8)
		b[0] = 1
		bgd.fs.dev.Write(b)

		bgd.Flags &= 0xFFFF ^ BG_INODE_UNINIT
		bgd.UpdateCsumAndWriteback()

		subInodeNum = 0
	} else {
		// Find free inode in bitmap
		for i := 0; i < int(bgd.fs.sb.InodePer_group/8); i++ {
			b := make([]byte, 1)
			bgd.fs.dev.Read(b)
			if b[0] != 0xFF {
				//log.Println("free at ", bgd.num, start, i)
				bitNum := bits.TrailingZeros8(^b[0])
				subInodeNum = int64(i)*8 + int64(bitNum)
				b[0] |= 1 << uint(bitNum)
				bgd.fs.dev.Seek(-1, 1)
				bgd.fs.dev.Write(b)
				break
			}
		}
	}

	if subInodeNum < 0 {
		//log.Println("!!!! bgd full !!!", bgd.num, bgd.Free_inodes_count_lo)
		return nil
	}

	if bgd.Flags & BG_INODE_ZEROED == 0 {
		bgd.fs.dev.Seek(bgd.GetInodeTableLoc() * bgd.fs.sb.GetBlockSize(), 0)
		bgd.fs.dev.Write(make([]byte, int64(bgd.fs.sb.InodePer_group) / int64(bgd.fs.sb.Inode_size)))
		bgd.Flags |= BG_INODE_ZEROED
		bgd.UpdateCsumAndWriteback()
	}

	// Update inode bitmap checksum
	checksummer := NewChecksummer(bgd.fs.sb)
	checksummer.Write(bgd.fs.sb.Uuid[:])
	bgd.fs.dev.Seek(start, 0)
	b := make([]byte, int64(bgd.fs.sb.InodePer_group) / 8)
	bgd.fs.dev.Read(b)
	checksummer.Write(b)
	bgd.Inode_bitmap_csum_lo = uint16(checksummer.Get() & 0xFFFF)
	bgd.Inode_bitmap_csum_hi = uint16(checksummer.Get() >> 16)

	bgd.Free_inodes_count_lo--
	bgd.Itable_unused_lo--
	bgd.UpdateCsumAndWriteback()

	bgd.fs.sb.Free_inodeCount--
	bgd.fs.sb.UpdateCsumAndWriteback()

	// Insert in Inode table
	inode := &Inode{
		Mode: 0,
		Links_count: 1,
		Flags: 524288, //TODO: what
		BlockOrExtents: [60]byte{0x0a, 0xf3, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00},
		fs: bgd.fs,
		address: bgd.GetInodeTableLoc() * bgd.fs.sb.GetBlockSize() + subInodeNum * int64(bgd.fs.sb.Inode_size),
		num: 1 + bgd.num * int64(bgd.fs.sb.InodePer_group) + subInodeNum,
	}
	inode.UpdateCsumAndWriteback()

	return inode
}

func (bgd *GroupDescriptor) setBitRange(offset int64, start int64, n int64) {
	if start < 0 || start > int64(bgd.fs.sb.BlockPer_group) {
		return
	}

	n = (n + bgd.fs.sb.GetBlockSize() - 1) / bgd.fs.sb.GetBlockSize()
	for i := start; i < start + n; i++ {
		b := make([]byte, 1)
		bgd.fs.dev.Seek(offset + i / 8, 0)
		bgd.fs.dev.Read(b)
		b[0] |= 1 << uint(i % 8)
		bgd.fs.dev.Seek(-1, 1)
		bgd.fs.dev.Write(b)
	}
	bgd.Free_blocks_count_lo-=uint16(n)
}

func (bgd *GroupDescriptor) GetFreeBlocks(n int64) (int64, int64) {
	// Find free block in bitmap
	start := bgd.GetBlockBitmapLoc() * bgd.fs.sb.GetBlockSize()

	if bgd.Flags & BG_BLOCK_UNINIT != 0 {
		bgd.fs.dev.Seek(start, 0)
		bgd.fs.dev.Write(make([]byte, bgd.fs.sb.BlockPer_group/8))
		bgd.Free_blocks_count_lo = uint16(bgd.fs.sb.BlockPer_group)

		if !bgd.fs.sb.FeatureRoCompatSparse_super() || bgd.num <= 1 || bgd.num % 3 == 0 || bgd.num % 5 == 0 || bgd.num % 7 == 0 {
			bgd.setBitRange(start, 0, bgd.fs.sb.GetBlockSize()+(bgd.fs.sb.numBlockGroups*32)+int64(bgd.fs.sb.Reserved_gdt_blocks)*bgd.fs.sb.GetBlockSize())
		}
		bgd.setBitRange(start, bgd.GetInodeBitmapLoc() - bgd.num * int64(bgd.fs.sb.BlockPer_group) * bgd.fs.sb.GetBlockSize(), int64(bgd.fs.sb.InodePer_group)/8)
		bgd.setBitRange(start, bgd.GetBlockBitmapLoc() - bgd.num * int64(bgd.fs.sb.BlockPer_group) * bgd.fs.sb.GetBlockSize(), int64(bgd.fs.sb.BlockPer_group)/8)
		bgd.setBitRange(start, bgd.GetInodeTableLoc() - bgd.num * int64(bgd.fs.sb.BlockPer_group) * bgd.fs.sb.GetBlockSize(), int64(bgd.fs.sb.Inode_size)*int64(bgd.fs.sb.InodePer_group)/8)
		bgd.UpdateCsumAndWriteback()

		blocksFree := int64(0)
		for i := int64(0); i < bgd.fs.sb.numBlockGroups; i++ {
			blocksFree += int64(bgd.fs.getBlockGroupDescriptor(i).Free_blocks_count_lo)
		}

		bgd.fs.sb.Free_blockCount_lo = uint32(blocksFree)
		bgd.fs.sb.UpdateCsumAndWriteback()

		bgd.Flags &= 0xFFFF ^ BG_BLOCK_UNINIT
		bgd.UpdateCsumAndWriteback()
	}

	subBlockNum := int64(-1)
	bgd.fs.dev.Seek(start, 0)
	for i := 0; i < int(bgd.fs.sb.BlockPer_group/8); i++ {
		b := make([]byte, 1)
		bgd.fs.dev.Read(b)
		if b[0] != 0xFF {
			bitNum := bits.TrailingZeros8(^b[0])
			numFree := bits.TrailingZeros8(uint8((uint16(b[0]) | 0x100) >> uint(bitNum)))
			//log.Println(bgd.num, i, b[0], bitNum, numFree, n)
			if n > int64(numFree) {
				n = int64(numFree)
			}
			subBlockNum = int64(i)*8 + int64(bitNum)
			b[0] |= (1 << uint(bitNum+int(n))) - 1
			//log.Println("Found free blocks. GD", bgd.num, subBlockNum, bitNum, b[0], n)
			bgd.fs.dev.Seek(-1, 1)
			bgd.fs.dev.Write(b)
			break
		}
	}

	if subBlockNum < 0 {
		return 0, 0
	}

	// Update block bitmap checksum
	checksummer := NewChecksummer(bgd.fs.sb)
	checksummer.Write(bgd.fs.sb.Uuid[:])
	bgd.fs.dev.Seek(start, 0)
	b := make([]byte, int64(bgd.fs.sb.ClusterPer_group) / 8)
	bgd.fs.dev.Read(b)
	checksummer.Write(b)
	bgd.Block_bitmap_csum_lo = uint16(checksummer.Get() & 0xFFFF)
	bgd.Block_bitmap_csum_hi = uint16(checksummer.Get() >> 16)

	newFreeBlocks := ((uint32(bgd.Free_blocks_count_hi) << 16) | uint32(bgd.Free_blocks_count_lo)) - uint32(n)
	bgd.Free_blocks_count_hi = uint16(newFreeBlocks >> 16)
	bgd.Free_blocks_count_lo = uint16(newFreeBlocks)
	bgd.UpdateCsumAndWriteback()

	bgd.fs.sb.Free_blockCount_lo-=uint32(n)
	bgd.fs.sb.UpdateCsumAndWriteback()

	return bgd.address / bgd.fs.sb.GetBlockSize() + subBlockNum - 1, n
}