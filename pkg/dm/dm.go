package dm

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/util"
	"log"
)

type blockDevice interface {
	Path() string
	activate() error
	active() bool
}

// Sectors is a data size unit for device mapper
// 1 sector = 512 bytes, Sectors basically divides via 512
// TODO: This needs to be a struct embedding an uint64
type Sectors uint64

func SectorsFromBytes(bytes interface{}) Sectors {
	switch bytes.(type) {
	case uint64:
		return Sectors(bytes.(uint64) / 512)
	case int64:
		return Sectors(bytes.(int64) / 512)
	case int:
		return Sectors(bytes.(int) / 512)
	}

	panic(fmt.Sprintf("invalid Sectors type: %T", bytes))
}

func (s Sectors) ToBytes() uint64 {
	return uint64(s * 512)
}

func dmsetup(args ...string) error {
	log.Printf("Running dmsetup: %q\n", args)
	_, err := util.ExecuteCommand("dmsetup", args...)
	return err
}

func resize2fs(device blockDevice) error {
	_, _ = util.ExecuteCommand("e2fsck", "-pf", device.Path())
	_, err := util.ExecuteCommand("resize2fs", device.Path())
	return err
}
