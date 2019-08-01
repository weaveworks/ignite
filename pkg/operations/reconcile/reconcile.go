package reconcile

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/storage/cache"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
	"github.com/weaveworks/ignite/pkg/util"
)

var c *client.Client

func ReconcileManifests(s *manifest.ManifestStorage) {
	startMetricsThread()

	// Wrap the Manifest Storage with a cache for better performance, and create a client
	c = client.NewClient(cache.NewCache(s))

	// These updates are coming from the SyncStorage
	for upd := range s.GetUpdateStream() {

		// Only care about VMs
		if upd.APIType.GetKind() != api.KindVM {
			log.Tracef("GitOps: Ignoring kind %s", upd.APIType.GetKind())
			continue
		}

		var vm *api.VM
		var err error
		if upd.Event == update.ObjectEventDelete {
			// As we know this VM was deleted, it wouldn't show up in a Get() call
			// Construct a temporary VM object for passing to the delete function
			vm = &api.VM{
				TypeMeta:   *upd.APIType.GetTypeMeta(),
				ObjectMeta: *upd.APIType.GetObjectMeta(),
				Status: api.VMStatus{
					State: api.VMStateStopped,
				},
			}
		} else {
			// Get the real API object
			vm, err = c.VMs().Get(upd.APIType.GetUID())
			if err != nil {
				log.Errorf("Getting VM %q returned an error: %v", upd.APIType.GetUID(), err)
				continue
			}

			// If the object was existent in the storage; validate it
			// Validate the VM object
			if err := validation.ValidateVM(vm).ToAggregate(); err != nil {
				log.Warn("Skipping %s of %s with UID %s, VM not valid %v.", upd.Event, upd.APIType.GetKind(), upd.APIType.GetUID(), err)
				continue
			}
		}

		// TODO: Paralellization
		switch upd.Event {
		case update.ObjectEventCreate:
			runHandle(func() error {
				return handleCreate(vm)
			})
		case update.ObjectEventModify:
			runHandle(func() error {
				return handleChange(vm)
			})
		case update.ObjectEventDelete:
			runHandle(func() error {
				// TODO: Temporary VM Object for removal
				return handleDelete(vm)
			})
		default:
			log.Infof("Unrecognized Git update type %s\n", upd.Event)
			continue
		}
	}
}

// TODO: Maybe paralellize these commands?
func runHandle(fn func() error) {
	if err := fn(); err != nil {
		log.Errorf("An error occurred when processing a VM update: %v\n", err)
	}
}

func handleCreate(vm *api.VM) error {
	var err error

	switch vm.Status.State {
	case api.VMStateCreated:
		err = create(vm)
	case api.VMStateRunning:
		err = start(vm)
	case api.VMStateStopped:
		log.Infof("VM %q was added to git with status Stopped, nothing to do\n", vm.GetUID())
	default:
		log.Infof("Unknown state of VM %q: %s", vm.GetUID().String(), vm.Status.State)
	}

	return err
}

func handleChange(vm *api.VM) error {
	var err error

	switch vm.Status.State {
	case api.VMStateCreated:
		err = fmt.Errorf("VM %q cannot changed into the Created state", vm.GetUID())
	case api.VMStateRunning:
		err = start(vm)
	case api.VMStateStopped:
		err = stop(vm)
	default:
		log.Infof("Unknown state of VM %q: %s", vm.GetUID().String(), vm.Status.State)
	}

	return err
}

func handleDelete(vm *api.VM) error {
	return remove(vm)
}

// TODO: Unify this with the "real" Create() method currently in cmd/
func create(vm *api.VM) error {
	log.Infof("Creating VM %q with name %q...", vm.GetUID(), vm.GetName())
	if err := ensureOCIImages(vm); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	return dmlegacy.AllocateAndPopulateOverlay(vm)
}

// ensureOCIImages imports the base/kernel OCI images if needed
func ensureOCIImages(vm *api.VM) error {
	// Check if a image with this name already exists, or import it
	image, err := operations.FindOrImportImage(c, vm.Spec.Image.OCIClaim.Ref)
	if err != nil {
		return err
	}

	// Populate relevant data from the Image on the VM object
	vm.SetImage(image)

	// Check if a kernel with this name already exists, or import it
	kernel, err := operations.FindOrImportKernel(c, vm.Spec.Kernel.OCIClaim.Ref)
	if err != nil {
		return err
	}

	// Populate relevant data from the Kernel on the VM object
	vm.SetKernel(kernel)

	// Save the file to disk. This will also write the file to /var/lib/firecracker for compability
	return c.VMs().Set(vm)
}

func start(vm *api.VM) error {
	// create the overlay if it doesn't exist
	if !util.FileExists(vm.OverlayFile()) {
		if err := create(vm); err != nil {
			return err
		}
	}

	log.Infof("Starting VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StartVM(vm, true)
}

func stop(vm *api.VM) error {
	log.Infof("Stopping VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StopVM(vm, true, false)
}

func remove(vm *api.VM) error {
	log.Infof("Removing VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.RemoveVM(c, vm)
}
