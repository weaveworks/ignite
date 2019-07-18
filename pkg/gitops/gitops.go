package gitops

import (
	"fmt"
	"sync"
	"time"

	"github.com/weaveworks/ignite/pkg/runtime/docker"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/gitops"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	vmMap map[meta.UID]*api.VM
	c     *client.Client
)

func RunLoop(url, branch string) error {
	log.Printf("Starting GitOps loop for repo at %q\n", url)
	log.Printf("Whenever changes are pushed to the %s branch, Ignite will apply the desired state locally\n", branch)
	log.Println("Initializing the Git repo...")

	s := gitops.NewGitOpsStorage(url, branch)
	// Wrap the GitOps storage with a cache for better performance
	c = client.NewClient(storage.NewCache(s))

	for {
		if !s.Ready() {
			// poll until the git repo is initialized
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Waiting for updates in the Git repo...")
		diff := s.WaitForUpdate()

		list, err := c.VMs().List()
		if err != nil {
			log.Warnf("Listing VMs returned an error: %v. Retrying...", err)
			continue
		}

		vmMap = mapVMs(list)

		wg := &sync.WaitGroup{}
		for _, file := range diff {
			// TODO: Construct a runVM object here and pass that instead of the "raw" API object
			vm := vmMap[file.APIType.UID]
			if vm == nil {
				if file.Type != gitops.UpdateTypeDeleted {
					// This is unexpected
					log.Warn("Skipping %s of %s with UID %s, no such object found through the client.", file.Type, file.APIType.GetKind(), file.APIType.GetUID())
					continue
				}

				// As we know this VM was deleted, it's logical that it wasn't found in the VMs().List() call above
				// Construct a temporary VM object for passing to the delete function
				vm = &api.VM{
					TypeMeta:   *file.APIType.TypeMeta,
					ObjectMeta: *file.APIType.ObjectMeta,
					Status: api.VMStatus{
						State: api.VMStateStopped,
					},
				}
			}

			// Construct the runtime object for this VM. This will also do defaulting
			// TODO: Consider cleanup like this?
			//defer metadata.Cleanup(runVM, false) // TODO: Handle silent
			//return metadata.Success(runVM)
			runVM := vmmd.WrapVM(vm)

			// TODO: At the moment there aren't running in parallel, shall they?
			switch file.Type {
			case gitops.UpdateTypeCreated:
				// TODO: Run this as a goroutine
				runHandle(wg, func() error {
					return handleCreate(runVM)
				})
			case gitops.UpdateTypeChanged:
				// TODO: Run this as a goroutine
				runHandle(wg, func() error {
					return handleChange(runVM)
				})
			case gitops.UpdateTypeDeleted:
				// TODO: Run this as a goroutine
				runHandle(wg, func() error {
					// TODO: Temporary VM Object for removal
					return handleDelete(runVM)
				})
			default:
				log.Printf("Unrecognized Git update type %s\n", file.Type)
				continue
			}
		}

		// wait for all goroutines to finish before the next sync round
		wg.Wait()
	}
}

func runHandle(wg *sync.WaitGroup, fn func() error) {
	wg.Add(1)
	defer wg.Done()

	if err := fn(); err != nil {
		log.Errorf("An error occurred when processing a VM update: %v\n", err)
	}
}

func mapVMs(vmlist []*api.VM) map[meta.UID]*api.VM {
	result := map[meta.UID]*api.VM{}
	for _, vm := range vmlist {
		result[vm.UID] = vm
	}

	return result
}

func handleCreate(vm *vmmd.VM) error {
	var err error

	switch vm.Status.State {
	case api.VMStateCreated:
		err = create(vm)
	case api.VMStateRunning:
		err = start(vm)
	case api.VMStateStopped:
		log.Printf("VM %q was added to git with status Stopped, nothing to do\n", vm.GetUID())
	default:
		log.Printf("Unknown state of VM %q: %s", vm.GetUID().String(), vm.Status.State)
	}

	return err
}

func handleChange(vm *vmmd.VM) error {
	var err error

	switch vm.Status.State {
	case api.VMStateCreated:
		err = fmt.Errorf("VM %q cannot changed into the Created state", vm.GetUID())
	case api.VMStateRunning:
		err = start(vm)
	case api.VMStateStopped:
		err = stop(vm)
	default:
		log.Printf("Unknown state of VM %q: %s", vm.GetUID().String(), vm.Status.State)
	}

	return err
}

func handleDelete(vm *vmmd.VM) error {
	return remove(vm)
}

// TODO: Unify this with the "real" Create() method currently in cmd/
func create(vm *vmmd.VM) error {
	log.Printf("Creating VM %q with name %q...", vm.GetUID(), vm.GetName())
	if err := ensureOCIImages(vm); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	return vm.AllocateAndPopulateOverlay()
}

// ensureOCIImages imports the base/kernel OCI images if needed
func ensureOCIImages(vm *vmmd.VM) error {
	if vm.Spec.Image.OCIClaim.Ref.IsUnset() {
		return fmt.Errorf("vm must specify image ref to run! image is empty for vm %s", vm.GetUID())
	}

	// Check if a image with this name already exists, or import it
	runImg, err := operations.FindOrImportImage(c, vm.Spec.Image.OCIClaim.Ref)
	if err != nil {
		return err
	}

	// Populate relevant data from the Image on the VM object
	vm.SetImage(runImg.Image)

	// Check if a kernel with this name already exists, or import it
	runKernel, err := operations.FindOrImportKernel(c, vm.Spec.Kernel.OCIClaim.Ref)
	if err != nil {
		return err
	}

	// Populate relevant data from the Kernel on the VM object
	vm.SetKernel(runKernel.Kernel)

	// Save the file to disk. This will also write the file to /var/lib/firecracker for compability
	return vm.Save()
}

func start(vm *vmmd.VM) error {
	// create the overlay if it doesn't exist
	if !util.FileExists(vm.OverlayFile()) {
		if err := create(vm); err != nil {
			return err
		}
	}

	log.Printf("Starting VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StartVM(vm, true)
}

func stop(vm *vmmd.VM) error {
	// Get the Docker client
	dc, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	log.Printf("Stopping VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StopVM(dc, vm, true, false)
}

func remove(vm *vmmd.VM) error {
	log.Printf("Removing VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.RemoveVM(c, vm)
}
