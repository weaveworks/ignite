package vmmd

//type blockDevice interface {
//	path() string
//}
//
//var _ blockDevice = &dmDevice{}
//var _ blockDevice = &dmPool{}
//var _ blockDevice = &loopDevice{}
//
//type dmDevice struct {
//	pool        *dmDevice
//	name        string
//	id          uint64
//	blocks      uint64
//	externalDev blockDevice
//}
//
//// blockSize specifies the data block size of the pool,
//// it should be between 128 (64KB) and 2097152 (1GB).
//// 128 is recommended if snapshotting a lot (like we do with layers).
//type dmPool struct {
//	dmDevice
//	metadataDev blockDevice
//	dataDev     blockDevice
//	blockSize   uint64
//}
//
//func newDMDevice(created bool, pool *dmDevice, name string, id, blocks uint64, externalDev blockDevice) (*dmDevice, error) {
//	// The volume is persistent in the metadata, it only needs to be generated once
//	if !created {
//		if err := dmsetup("message", pool.path(), "0", fmt.Sprintf("create_thin %d", id)); err != nil {
//			return nil, err
//		}
//	}
//
//	newDev := &dmDevice{
//		pool:        pool,
//		name:        name,
//		id:          id,
//		blocks:      blocks,
//		externalDev: externalDev,
//	}
//
//	return newDev.activate()
//}
//
//func (d *dmDevice) createSnapshot(created bool, name string) (*dmDevice, error) {
//	snapshotID := d.id + 1
//
//	// The snapshot is persistent in the metadata, it only needs to be generated once
//	if !created {
//		if err := dmsetup("suspend", d.path()); err != nil {
//			return nil, err
//		}
//
//		if err := dmsetup("message", d.pool.path(), "0",
//			fmt.Sprintf("create_snap %d %d", snapshotID, d.id)); err != nil {
//			return nil, err
//		}
//
//		if err := dmsetup("resume", d.path()); err != nil {
//			return nil, err
//		}
//	}
//
//	snapshotDev := &dmDevice{
//		pool:        d.pool,
//		name:        name,
//		id:          snapshotID,
//		blocks:      d.blocks,
//		externalDev: d.externalDev,
//	}
//
//	return snapshotDev.activate()
//}
//
//func newDMPool(name string, blocks, blockSize uint64, metadataDev, dataDev blockDevice) (*dmPool, error) {
//	pool := &dmPool{
//		dmDevice: dmDevice{
//			name:   name,
//			blocks: blocks,
//		},
//		metadataDev: metadataDev,
//		dataDev:     dataDev,
//		blockSize:   blockSize,
//	}
//
//	return pool.activate()
//}
//
//func (d *dmDevice) activate() (*dmDevice, error) {
//	dmTable := fmt.Sprintf("0 %d thin %s %d",
//		d.blocks,
//		d.pool.path(),
//		d.id,
//	)
//
//	if d.externalDev != nil {
//		dmTable = fmt.Sprintf("%s %s", dmTable, d.externalDev.path())
//	}
//
//	if err := dmsetup("create", d.name, "--table", dmTable); err != nil {
//		return nil, err
//	}
//
//	return d, nil
//}
//
//func (d *dmPool) activate() (*dmPool, error) {
//	dmTable := fmt.Sprintf("0 %d thin-pool %s %s %d 0",
//		d.blocks,
//		d.metadataDev.path(),
//		d.dataDev.path(),
//		d.blockSize,
//	)
//
//	if err := dmsetup("create", d.name, "--table", dmTable); err != nil {
//		return nil, err
//	}
//
//	return d, nil
//}
//
//func (d *dmDevice) path() string {
//	return path.Join("/dev/mapper", d.name)
//}
//
//type devNames []string
//
//// Order matters for removal
//func (md *VMMetadata) newDevNames() devNames {
//	return []string{
//		constants.IGNITE_PREFIX + md.ID.String(),
//		constants.IGNITE_PREFIX + "base-" + md.ID.String(),
//		constants.IGNITE_PREFIX + "pool-" + md.ID.String(),
//	}
//}
//
//func (d devNames) overlay() string {
//	return d[0]
//}
//
//func (d devNames) base() string {
//	return d[1]
//}
//
//func (d devNames) pool() string {
//	return d[2]
//}
//
//func (d devNames) all() []string {
//	return d
//}
//
//func (md *VMMetadata) NewVMOverlay() error {
//	names := md.newDevNames()
//
//	// Return if the overlay is already setup
//	// TODO: Check this individually for each volume
//	if util.FileExists((&dmDevice{name: names.overlay()}).path()) {
//		return nil
//	}
//
//	// Setup loop device for the metadata
//	metadataDev, err := newLoopDevice(path.Join(md.ObjectPath(), constants.VM_METADATA_FILE), false)
//	if err != nil {
//		return err
//	}
//
//	// Setup loop device for the data
//	dataDev, err := newLoopDevice(path.Join(md.ObjectPath(), constants.VM_DATA_FILE), false)
//	if err != nil {
//		return err
//	}
//
//	// The pool size should be the size of the data device
//	poolSize, err := dataDev.Size512K()
//	if err != nil {
//		return err
//	}
//
//	// Create the thin provisioning pool
//	pool, err := newDMPool(names.pool(), poolSize, 128, metadataDev, dataDev)
//	if err != nil {
//		return err
//	}
//
//	// Setup loop device for the image
//	imageDev, err := newLoopDevice(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID.String(), constants.IMAGE_FS), true)
//	if err != nil {
//		return err
//	}
//
//	// Detect if we're running for the first time, this is needed for triggering volume/snapshot creation
//	created := md.VMOD().VolumesCreated
//
//	// Create the base device, which is an external snapshot of the image
//	baseDev, err := newDMDevice(created, &pool.dmDevice, names.base(), 0, pool.blocks, imageDev)
//	if err != nil {
//		return err
//	}
//
//	// Resize the filesystem to fill the base
//	if err := resize2fs(baseDev); err != nil {
//		return err
//	}
//
//	// TODO: Save/return this overlay device
//	if _, err = baseDev.createSnapshot(created, names.overlay()); err != nil {
//		return err
//	}
//
//	// Mark the volumes to be created
//	if !created {
//		md.VMOD().VolumesCreated = true
//		if err := md.Save(); err != nil {
//			return err
//		}
//	}
//
//	// By detaching the loop devices after setting up thin provisioning
//	// they get automatically removed when the thin volumes/snapshots are removed.
//	if err := metadataDev.Detach(); err != nil {
//		return err
//	}
//
//	if err := dataDev.Detach(); err != nil {
//		return err
//	}
//
//	if err := imageDev.Detach(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (md *VMMetadata) OverlayDev() string {
//	return (&dmDevice{name: md.newDevNames().overlay()}).path()
//}
//
//func (md *VMMetadata) RemoveOverlay() error {
//	devNames := md.newDevNames().all()
//	dmArgs := append(make([]string, 0, len(devNames)+1), "remove")
//	return dmsetup(append(dmArgs, devNames...)...)
//}
//
//type loopDevice struct {
//	losetup.Device
//}
//
//func newLoopDevice(file string, readOnly bool) (*loopDevice, error) {
//	dev, err := losetup.Attach(file, 0, readOnly)
//	if err != nil {
//		return nil, fmt.Errorf("failed to setup loop device for %q: %v", file, err)
//	}
//
//	return &loopDevice{dev}, nil
//}
//
//func (ld *loopDevice) Size512K() (uint64, error) {
//	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
//	if err != nil {
//		return 0, err
//	}
//
//	// Remove the trailing newline and parse to uint64
//	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
//}
//
//func (ld *loopDevice) path() string {
//	return ld.Path()
//}
//
//func dmsetup(args ...string) error {
//	log.Printf("Running dmsetup: %q\n", args)
//	_, err := util.ExecuteCommand("dmsetup", args...)
//	return err
//}
//
//func resize2fs(device blockDevice) error {
//	_, _ = util.ExecuteCommand("e2fsck", "-pf", device.path())
//	_, err := util.ExecuteCommand("resize2fs", device.path())
//	return err
//}
