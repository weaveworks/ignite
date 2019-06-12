package vmmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strconv"

	"github.com/freddierice/go-losetup"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	poolPrefix = constants.IGNITE_PREFIX + "pool-"
	basePrefix = constants.IGNITE_PREFIX + "base-"
)

type VMFilesystem struct{}

type blockDev interface {
	path() string
}

var _ blockDev = &dmDevice{}
var _ blockDev = &dmPool{}
var _ blockDev = &loopDevice{}

type dmDevice struct {
	pool        *dmDevice
	name        string
	id          uint64
	blocks      uint64
	externalDev blockDev
}

// blockSize specifies the data block size of the pool,
// it should be between 128 (64KB) and 2097152 (1GB).
// 128 is recommended if snapshotting a lot (like we do with layers).
type dmPool struct {
	dmDevice
	metadataDev blockDev
	dataDev     blockDev
	blockSize   uint64
}

func newDMDevice(pool *dmDevice, name string, id, blocks uint64, externalDev blockDev) (*dmDevice, error) {
	// The volume needs to be generated, but persists across rebuilds, ignore the error for now
	// TODO: Better error handling
	_ = dmsetup("message", pool.path(), "0", fmt.Sprintf("create_thin %d", id))

	return activateDMDevice(pool, name, id, blocks, externalDev)
}

func activateDMDevice(pool *dmDevice, name string, id, blocks uint64, externalDev blockDev) (*dmDevice, error) {
	dev := &dmDevice{
		pool:        pool,
		name:        name,
		id:          id,
		blocks:      blocks,
		externalDev: externalDev,
	}

	if err := dev.create(); err != nil {
		return nil, err
	}

	return dev, nil
}

func newDMPool(name string, blocks, blockSize uint64, metadataDev, dataDev blockDev) (*dmPool, error) {
	pool := &dmPool{
		dmDevice: dmDevice{
			name:   name,
			blocks: blocks,
		},
		metadataDev: metadataDev,
		dataDev:     dataDev,
		blockSize:   blockSize,
	}

	if err := pool.create(); err != nil {
		return nil, err
	}

	return pool, nil
}

func (d *dmDevice) create() error {
	dmTable := fmt.Sprintf("0 %d thin %s %d",
		d.blocks,
		d.pool.path(),
		d.id,
	)

	if d.externalDev != nil {
		dmTable = fmt.Sprintf("%s %s", dmTable, d.externalDev.path())
	}

	if err := dmsetup("create", d.name, "--table", dmTable); err != nil {
		return err
	}

	return nil
}

func (d *dmPool) create() error {
	dmTable := fmt.Sprintf("0 %d thin-pool %s %s %d 0",
		d.blocks,
		d.metadataDev.path(),
		d.dataDev.path(),
		d.blockSize,
	)

	if err := dmsetup("create", d.name, "--table", dmTable); err != nil {
		return err
	}

	return nil
}

func (d *dmDevice) path() string {
	return path.Join("/dev/mapper", d.name)
}

func (md *VMMetadata) NewVMOverlay() error {
	poolName := poolPrefix + md.ID.String()
	baseDevName := basePrefix + md.ID.String()
	overlayDevName := constants.IGNITE_PREFIX + md.ID.String()

	// Return if the overlay is already setup
	if util.FileExists((&dmDevice{name: overlayDevName}).path()) {
		return nil
	}

	// Setup loop device for the metadata
	metadataDev, err := newLoopDev(path.Join(md.ObjectPath(), constants.METADATA_FILE), false)
	if err != nil {
		return err
	}

	// Setup loop device for the data
	dataDev, err := newLoopDev(path.Join(md.ObjectPath(), constants.DATA_FILE), false)
	if err != nil {
		return err
	}

	poolSize, err := dataDev.Size512K()
	if err != nil {
		return err
	}

	pool, err := newDMPool(poolName, poolSize, 128, metadataDev, dataDev)
	if err != nil {
		return err
	}

	// Setup loop device for the image
	imageDev, err := newLoopDev(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID.String(), constants.IMAGE_FS), false)
	if err != nil {
		return err
	}

	baseDev, err := newDMDevice(&pool.dmDevice, baseDevName, 0, pool.blocks, imageDev)
	if err != nil {
		return err
	}

	// Resize the filesystem to fill the base
	if err := resize2fs(baseDev); err != nil {
		return err
	}

	// TODO: Save/return this overlay device
	if _, err = baseDev.createSnapshot(overlayDevName); err != nil {
		return err
	}

	return nil
}

func (md *VMMetadata) OverlayDev() string {
	return path.Join("/dev/mapper", constants.IGNITE_PREFIX+md.ID.String())
}

func (md *VMMetadata) RemoveOverlay() error {
	log.Println("Overlay remove: stub")
	return nil
}

//func (md *VMMetadata) SnapshotDev() string {
//	return path.Join("/dev/mapper", constants.IGNITE_PREFIX+md.ID.String())
//}
//
//func (md *VMMetadata) SetupSnapshot() error {
//	device := constants.IGNITE_PREFIX + md.ID.String()
//	devicePath := md.SnapshotDev()
//
//	// Return if the snapshot is already setup
//	if util.FileExists(devicePath) {
//		return nil
//	}
//
//	// Setup loop device for the image
//	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID.String(), constants.IMAGE_FS), true)
//	if err != nil {
//		return err
//	}
//
//	// Setup loop device for the VM overlay
//	overlayLoop, err := newLoopDev(path.Join(md.ObjectPath(), constants.OVERLAY_FILE), false)
//	if err != nil {
//		return err
//	}
//
//	imageLoopSize, err := imageLoop.Size512K()
//	if err != nil {
//		return err
//	}
//
//	overlayLoopSize, err := overlayLoop.Size512K()
//	if err != nil {
//		return err
//	}
//
//	// If the overlay is larger than the base image, we need to set up an additional dm device
//	// which will contain the image and additional zero space (which reads zeros and discards writes).
//	// This is fine, because all writes will target the overlay snapshot and not the read-only image.
//	// The newly generated larger device will then be used for creating the snapshot (which is always
//	// as large as the device backing it).
//
//	basePath := imageLoop.Path()
//	if overlayLoopSize > imageLoopSize {
//		// "0 8388608 linear /dev/loop0 0"
//		// "8388608 12582912 zero"
//		dmBaseTable := []byte(fmt.Sprintf("0 %d linear %s 0\n%d %d zero", imageLoopSize, imageLoop.Path(), imageLoopSize, overlayLoopSize))
//
//		baseDevice := fmt.Sprintf("%s-base", device)
//		if err := runDMSetup(baseDevice, dmBaseTable); err != nil {
//			return err
//		}
//
//		basePath = fmt.Sprintf("/dev/mapper/%s", baseDevice)
//	}
//
//	// "0 8388608 snapshot /dev/{loop0,mapper/ignite-<uid>-base} /dev/loop1 P 8"
//	dmTable := []byte(fmt.Sprintf("0 %d snapshot %s %s P 8", overlayLoopSize, basePath, overlayLoop.Path()))
//
//	if err := runDMSetup(device, dmTable); err != nil {
//		return err
//	}
//
//	// Repair the filesystem in case it has errors
//	// e2fsck throws an error if the filesystem gets repaired, so just ignore it
//	_, _ = util.ExecuteCommand("e2fsck", "-p", "-f", devicePath)
//
//	// If the overlay is larger than the image, call resize2fs to make the filesystem fill the overlay
//	if overlayLoopSize > imageLoopSize {
//		if _, err := util.ExecuteCommand("resize2fs", devicePath); err != nil {
//			return err
//		}
//	}
//
//	// By detaching the loop devices after setting up the snapshot
//	// they get automatically removed when the snapshot is removed.
//	if err := imageLoop.Detach(); err != nil {
//		return err
//	}
//
//	if err := overlayLoop.Detach(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (md *VMMetadata) RemoveSnapshot() error {
//	dmArgs := []string{
//		"remove",
//		md.SnapshotDev(),
//	}
//
//	// If the base device is visible in "dmsetup", we should remove it
//	// The device itself is not forwarded to docker, so we can't query its path
//	// TODO: Improve this detection
//	baseDev := fmt.Sprintf("%s-base", constants.IGNITE_PREFIX+md.ID.String())
//	if _, err := util.ExecuteCommand("dmsetup", "info", baseDev); err == nil {
//		dmArgs = append(dmArgs, baseDev)
//	}
//
//	if _, err := util.ExecuteCommand("dmsetup", dmArgs...); err != nil {
//		return err
//	}
//
//	return nil
//}

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

func (ld *loopDevice) path() string {
	return ld.Path()
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

func (d *dmDevice) createSnapshot(name string) (*dmDevice, error) {
	snapshotID := d.id + 1

	if err := dmsetup("suspend", d.path()); err != nil {
		return nil, err
	}

	// The snapshot needs to be generated, but persists across rebuilds, ignore the error for now
	// TODO: Better error handling
	if err := dmsetup("message", d.pool.path(), "0",
		fmt.Sprintf("create_snap %d %d", snapshotID, d.id)); err != nil {
		//return nil, err
	}

	if err := dmsetup("resume", d.path()); err != nil {
		return nil, err
	}

	return activateDMDevice(d.pool, name, snapshotID, d.blocks, d.externalDev)
}

func dmsetup(args ...string) error {
	log.Printf("Running dmsetup: %q\n", args)
	_, err := util.ExecuteCommand("dmsetup", args...)
	return err
}

func resize2fs(device blockDev) error {
	_, _ = util.ExecuteCommand("e2fsck", "-pf", device.path())
	_, err := util.ExecuteCommand("resize2fs", device.path())
	return err
}
