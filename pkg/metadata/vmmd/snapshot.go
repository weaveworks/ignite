package vmmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"

	"github.com/freddierice/go-losetup"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

func (md *VM) SnapshotDev() string {
	return path.Join("/dev/mapper", constants.IGNITE_PREFIX+md.GetUID().String())
}

func (md *VM) SetupSnapshot() error {
	device := constants.IGNITE_PREFIX + md.GetUID().String()
	devicePath := md.SnapshotDev()

	// Return if the snapshot is already setup
	if util.FileExists(devicePath) {
		return nil
	}

	// Setup loop device for the image
	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, md.Spec.Image.UID.String(), constants.IMAGE_FS), true)
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

	// If the overlay is larger than the base image, we need to set up an additional dm device
	// which will contain the image and additional zero space (which reads zeros and discards writes).
	// This is fine, because all writes will target the overlay snapshot and not the read-only image.
	// The newly generated larger device will then be used for creating the snapshot (which is always
	// as large as the device backing it).

	basePath := imageLoop.Path()
	if overlayLoopSize > imageLoopSize {
		// "0 8388608 linear /dev/loop0 0"
		// "8388608 12582912 zero"
		dmBaseTable := []byte(fmt.Sprintf("0 %d linear %s 0\n%d %d zero", imageLoopSize, imageLoop.Path(), imageLoopSize, overlayLoopSize))

		baseDevice := fmt.Sprintf("%s-base", device)
		if err := runDMSetup(baseDevice, dmBaseTable); err != nil {
			return err
		}

		basePath = fmt.Sprintf("/dev/mapper/%s", baseDevice)
	}

	// "0 8388608 snapshot /dev/{loop0,mapper/ignite-<uid>-base} /dev/loop1 P 8"
	dmTable := []byte(fmt.Sprintf("0 %d snapshot %s %s P 8", overlayLoopSize, basePath, overlayLoop.Path()))

	if err := runDMSetup(device, dmTable); err != nil {
		return err
	}

	// Repair the filesystem in case it has errors
	// e2fsck throws an error if the filesystem gets repaired, so just ignore it
	_, _ = util.ExecuteCommand("e2fsck", "-p", "-f", devicePath)

	// If the overlay is larger than the image, call resize2fs to make the filesystem fill the overlay
	if overlayLoopSize > imageLoopSize {
		if _, err := util.ExecuteCommand("resize2fs", devicePath); err != nil {
			return err
		}
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

func (md *VM) RemoveSnapshot() error {
	dmArgs := []string{
		"remove",
		md.SnapshotDev(),
	}

	// If the base device is visible in "dmsetup", we should remove it
	// The device itself is not forwarded to docker, so we can't query its path
	// TODO: Improve this detection
	baseDev := fmt.Sprintf("%s-base", constants.IGNITE_PREFIX+md.GetUID())
	if _, err := util.ExecuteCommand("dmsetup", "info", baseDev); err == nil {
		dmArgs = append(dmArgs, baseDev)
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

// dmsetup uses stdin to read multiline tables, this is a helper function for that
func runDMSetup(name string, table []byte) error {
	cmd := exec.Command("dmsetup", "create", name)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, err := stdin.Write(table); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command %q exited with %q: %v", cmd.Args, out, err)
	}

	return nil
}
