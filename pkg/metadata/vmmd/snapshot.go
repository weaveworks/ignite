package vmmd

import (
	"fmt"
	"github.com/freddierice/go-losetup"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

func (md *VMMetadata) SnapshotDev() string {
	return path.Join("/dev/mapper", constants.IGNITE_PREFIX+md.ID)
}

func (md *VMMetadata) SetupSnapshot() error {
	device := constants.IGNITE_PREFIX + md.ID
	devicePath := md.SnapshotDev()

	// Return if the snapshot is already setup
	if util.FileExists(devicePath) {
		return nil
	}

	// Setup loop device for the image
	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID, constants.IMAGE_FS), true)
	if err != nil {
		return err
	}

	// Setup loop device for the VM overlay
	overlayLoop, err := newLoopDev(path.Join(md.ObjectPath(), constants.OVERLAY_FILE), false)
	if err != nil {
		return err
	}

	imageLoopSize, err := imageLoop.Size512K()
	if err != nil {
		return err
	}

	overlayLoopSize, err := overlayLoop.Size512K()
	if err != nil {
		return err
	}

	basePath := imageLoop.Path()

	// TODO: Fix support for overlays smaller than the image
	if overlayLoopSize > imageLoopSize {
		fmt.Println("Overlay larger than image!")

		dmBaseTable := []byte(fmt.Sprintf("0 %d linear %s 0\n%d %d zero", imageLoopSize, imageLoop.Path(), imageLoopSize, overlayLoopSize))

		dmBaseArgs := []string{
			"create",
			device + "-base",
		}

		if _, err := util.ExecuteCommandStdin("dmsetup", dmBaseTable, dmBaseArgs...); err != nil {
			return err
		}

		basePath = fmt.Sprintf("/dev/mapper/%s-base", device)
		fmt.Println("Success!")
	}

	// dmsetup create newdev --table "0 8388608 snapshot /dev/loop0 /dev/loop1 P 8"
	dmTable := []string{
		"0",
		strconv.FormatUint(overlayLoopSize, 10),
		"snapshot",
		basePath,
		overlayLoop.Path(),
		"P",
		"8",
	}

	dmArgs := []string{
		"create",
		device,
		"--table",
		strings.Join(dmTable, " "),
	}

	if _, err := util.ExecuteCommand("dmsetup", dmArgs...); err != nil {
		return err
	}

	// Call resize2fs to make the filesystem fill the overlay
	if _, err := util.ExecuteCommand("resize2fs", devicePath); err != nil {
		return err
	}

	// By detaching the loop devices after setting up the snapshot
	// they get automatically removed when the snapshot is removed.
	if err := imageLoop.Detach(); err != nil {
		return err
	}

	if err := overlayLoop.Detach(); err != nil {
		return err
	}

	return nil
}

func (md *VMMetadata) RemoveSnapshot() error {
	dmArgs := []string{
		"remove",
		md.SnapshotDev(),
	}

	if _, err := util.ExecuteCommand("dmsetup", dmArgs...); err != nil {
		return err
	}

	return nil
}

type loopDevice struct {
	losetup.Device
}

func newLoopDev(file string, readOnly bool) (*loopDevice, error) {
	dev, err := losetup.Attach(file, 0, readOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to setup loop device for %q: %v", file, err)
	}

	return &loopDevice{dev}, nil
}

func (ld *loopDevice) Size512K() (uint64, error) {
	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
	if err != nil {
		return 0, err
	}

	// Remove the trailing newline and parse to uint64
	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
}
