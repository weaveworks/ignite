package gitops

import (
	"fmt"
	"log"
	"sync"
	"time"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/storage/gitops"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	vmMap map[meta.UID]*api.VM
	s     *gitops.GitOpsStorage
	c     *client.Client
)

func RunLoop(url, branch string) error {
	log.Printf("Starting GitOps loop for repo at %q\n", url)
	log.Printf("Whenever changes are pushed the %s branch, Ignite will apply the desired state locally\n", branch)
	log.Println("Initializing the Git repo...")

	s = gitops.NewGitOpsStorage(url, branch)
	c = client.NewClient(s)

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
			fmt.Printf("list err %v", err)
			continue
		}
		vmMap = mapVMs(list)

		wg := &sync.WaitGroup{}
		for _, file := range diff {
			vm := vmMap[file.APIType.UID]

			// TODO: At the moment there aren't running in parallel, shall they?
			switch file.Type {
			case gitops.UpdateTypeCreated:
				go runHandle(wg, func() error {
					return handleCreate(vm)
				})
			case gitops.UpdateTypeChanged:
				go runHandle(wg, func() error {
					return handleChange(vm)
				})
			case gitops.UpdateTypeDeleted:
				go runHandle(wg, func() error {
					return handleDelete(file.APIType)
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
		log.Printf("[WARNING] An error occurred when processing a VM update: %v\n", err)
	}
}

func mapVMs(vmlist []*api.VM) map[meta.UID]*api.VM {
	result := map[meta.UID]*api.VM{}
	for _, vm := range vmlist {
		result[vm.UID] = vm
	}
	return result
}

func handleCreate(vm *api.VM) error {
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
		log.Printf("Unknown state of VM %q: %s", vm.GetUID().String(), vm.Status.State)
	}
	return err
}

func handleDelete(obj *meta.APIType) error {
	return remove(obj)
}

// TODO: use a real filter here when ready
func findImageByName(list []*api.Image, name string) *api.Image {
	for _, obj := range list {
		if obj.GetName() == name {
			return obj
		}
	}
	return nil
}
func findKernelByName(list []*api.Kernel, name string) *api.Kernel {
	for _, obj := range list {
		if obj.GetName() == name {
			return obj
		}
	}
	return nil
}

// TODO: Unify this with the "real" Create() method currently in cmd/
func create(vm *api.VM) error {
	log.Printf("Creating VM %q with name %q...", vm.GetUID(), vm.GetName())
	runVM, err := prepareVM(vm)
	if err != nil {
		return err
	}

	// Save the file to disk. This will also write the file to /var/lib/firecracker for compability
	if err := runVM.Save(); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	return runVM.AllocateAndPopulateOverlay()
}

// prepareVM takes a VM API object, finds/populates its dependencies (image/kernel) and finally
// returns the runtime VM object
func prepareVM(vm *api.VM) (*vmmd.VM, error) {
	if vm.Spec.Image.OCIClaim == nil {
		return nil, fmt.Errorf("vm must specify image to run! image is empty for vm %s", vm.GetUID())
	}

	imgName := vm.Spec.Image.OCIClaim.Ref
	imgName, _ = metadata.NewNameWithLatest(imgName, meta.KindImage)
	imgs, err := c.Images().List()
	if err != nil {
		return nil, err
	}

	var runImg *imgmd.Image
	img := findImageByName(imgs, imgName)
	if img == nil {
		if runImg, err = operations.ImportImage(imgName); err != nil {
			return nil, fmt.Errorf("failed to import image %s %v", imgName, err)
		}
	} else {
		if runImg, err = imgmd.NewImage(img.UID, &img.Name, img); err != nil {
			return nil, err
		}
	}

	// check if a kernel with this name already exists, or import it
	kernels, err := c.Kernels().List()
	if err != nil {
		return nil, err
	}

	var runKernel *kernmd.Kernel
	kernel := findKernelByName(kernels, imgName)
	if kernel == nil {
		if runKernel, err = operations.ImportKernelFromImage(runImg); err != nil {
			return nil, fmt.Errorf("failed to import kernel %s %v", imgName, err)
		}
	} else {
		if runKernel, err = kernmd.NewKernel(kernel.UID, &kernel.Name, kernel); err != nil {
			return nil, err
		}
	}

	// populate the image/kernel ID fields to use when running the VM
	vm.Status.Image.UID = runImg.UID
	vm.Status.Kernel.UID = runKernel.UID

	// Create new metadata for the VM
	return vmmd.NewVM(vm.ObjectMeta.UID, &vm.ObjectMeta.Name, vm)
	// TODO: Consider cleanup like this?
	//defer metadata.Cleanup(runVM, false) // TODO: Handle silent
	//return metadata.Success(runVM)
}

func start(vm *api.VM) error {
	runVM, err := prepareVM(vm)
	if err != nil {
		return err
	}
	// create the overlay if it doesn't exist
	if !util.FileExists(runVM.OverlayFile()) {
		if err := create(vm); err != nil {
			return err
		}
	}

	log.Printf("Starting VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StartVM(runVM, true)
}

func stop(vm *api.VM) error {
	log.Printf("Stopping VM %q with name %q...", vm.GetUID(), vm.GetName())
	return operations.StopVM(vm, true, false)
}

func remove(obj *meta.APIType) error {
	log.Printf("Removing VM %q with name %q...", obj.GetUID(), obj.GetName())
	return operations.RemoveVM(c, obj)
}
