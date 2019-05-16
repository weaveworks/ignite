package gexto

import (
	"hash/crc32"
	"encoding/binary"
	"io"
)

type Checksummer interface {
	io.Writer
	WriteUint32(uint32)
	Get() uint32
}

func NewChecksummer(sb *Superblock) Checksummer {
	return &checksummer{
		sb: sb,
		val: 0,
		table: crc32.MakeTable(crc32.Castagnoli), // TODO: Check crc used in sb?
	}
}

type checksummer struct {
	sb          *Superblock
	val         uint32
	table       *crc32.Table
}

func (cs *checksummer) Write(b []byte) (n int, err error) {
	cs.val = crc32.Update(cs.val, cs.table, b)
	return len(b), nil
}

func (cs *checksummer) WriteUint32(x uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, x)
	cs.Write(b)
}

func (cs *checksummer) Get() uint32 {
	return ^cs.val
}