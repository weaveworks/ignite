package losetup

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Info is a datastructure that holds relevant information about a file backed
// loopback device.
type Info struct {
	Device         uint64
	INode          uint64
	RDevice        uint64
	Offset         uint64
	SizeLimit      uint64
	Number         uint32
	EncryptType    uint32
	EncryptKeySize uint32
	Flags          uint32
	FileName       [NameSize]byte
	CryptName      [NameSize]byte
	EncryptKey     [KeySize]byte
	Init           [2]uint64
}

// GetInfo returns information about a loop device
func (device Device) GetInfo() (Info, error) {
	f, err := device.open()
	if err != nil {
		return Info{}, fmt.Errorf("could not open %v: %v", device, err)
	}
	defer f.Close()

	return getInfo(f.Fd())
}

func getInfo(fd uintptr) (Info, error) {
	retInfo := Info{}
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, GetStatus64, uintptr(unsafe.Pointer(&retInfo)))
	if errno == unix.ENXIO {
		return Info{}, fmt.Errorf("device not backed by a file")
	} else if errno != 0 {
		return Info{}, fmt.Errorf("could not get info about %v (err: %d): %v", errno, errno)
	}

	return retInfo, nil
}

// SetInfo sets options in the loop device.
func (device Device) SetInfo(info Info) error {
	f, err := device.open()
	if err != nil {
		return fmt.Errorf("could not open %v: %v", device, err)
	}
	defer f.Close()

	return setInfo(f.Fd(), info)
}

func setInfo(fd uintptr, info Info) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, SetStatus64, uintptr(unsafe.Pointer(&info)))
	if errno == unix.ENXIO {
		return fmt.Errorf("device not backed by a file")
	} else if errno != 0 {
		return fmt.Errorf("could not get info about %v (err: %d): %v", errno, errno)
	}

	return nil
}
