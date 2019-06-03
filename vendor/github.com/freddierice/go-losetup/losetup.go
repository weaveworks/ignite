package losetup

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Add will add a loopback device if it does not exist already.
func (device Device) Add() error {
	ctrl, err := os.OpenFile(LoopControlPath, os.O_RDWR, 0660)
	if err != nil {
		return fmt.Errorf("could not open %v: %v", LoopControlPath, err)
	}
	defer ctrl.Close()
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, ctrl.Fd(), CtlAdd, uintptr(device.number))
	if errno == unix.EEXIST {
		return fmt.Errorf("device already exits")
	}
	if errno != 0 {
		return fmt.Errorf("could not add device (err: %d): %v", errno, errno)
	}
	return nil
}

// Remove will remove a loopback device if it is not busy.
func (device Device) Remove() error {
	ctrl, err := os.OpenFile(LoopControlPath, os.O_RDWR, 0660)
	if err != nil {
		return fmt.Errorf("could not open %v: %v", LoopControlPath, err)
	}
	defer ctrl.Close()
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, ctrl.Fd(), CtlRemove, uintptr(device.number))
	if errno == unix.EBUSY {
		return fmt.Errorf("could not remove, device in use")
	}
	if errno != 0 {
		return fmt.Errorf("could not remove (err: %d): %v", errno, errno)
	}
	return nil
}

// GetFree searches for the first free loopback device. If it cannot find one,
// it will attempt to create one. If anything fails, GetFree will return an
// error.
func GetFree() (Device, error) {
	ctrl, err := os.OpenFile(LoopControlPath, os.O_RDWR, 0660)
	if err != nil {
		return Device{}, fmt.Errorf("could not open %v: %v", LoopControlPath, err)
	}
	defer ctrl.Close()
	dev, _, errno := unix.Syscall(unix.SYS_IOCTL, ctrl.Fd(), CtlGetFree, 0)
	if dev < 0 {
		return Device{}, fmt.Errorf("could not get free device (err: %d): %v", errno, errno)
	}
	return Device{number: uint64(dev), flags: os.O_RDWR}, nil
}

// Attach attaches backingFile to the loopback device starting at offset. If ro
// is true, then the file is attached read only.
func Attach(backingFile string, offset uint64, ro bool) (Device, error) {
	var dev Device

	flags := os.O_RDWR
	if ro {
		flags = os.O_RDONLY
	}

	back, err := os.OpenFile(backingFile, flags, 0660)
	if err != nil {
		return dev, fmt.Errorf("could not open backing file: %v", err)
	}
	defer back.Close()

	dev, err = GetFree()
	if err != nil {
		return dev, err
	}
	dev.flags = flags

	loopFile, err := dev.open()
	if err != nil {
		return dev, fmt.Errorf("could not open loop device: %v", err)
	}
	defer loopFile.Close()

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, loopFile.Fd(), SetFd, back.Fd())
	if errno == 0 {
		info := Info{}
		copy(info.FileName[:], []byte(backingFile))
		info.Offset = offset
		if err := setInfo(loopFile.Fd(), info); err != nil {
			unix.Syscall(unix.SYS_IOCTL, loopFile.Fd(), ClrFd, 0)
			return dev, fmt.Errorf("could not set info")
		}
	}

	return dev, nil
}

// Detach removes the file backing the device.
func (device Device) Detach() error {

	loopFile, err := os.OpenFile(device.Path(), os.O_RDONLY, 0660)
	if err != nil {
		return fmt.Errorf("could not open loop device")
	}
	defer loopFile.Close()

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, loopFile.Fd(), ClrFd, 0)
	if errno != 0 {
		return fmt.Errorf("error clearing loopfile: %v", errno)
	}

	return nil
}
