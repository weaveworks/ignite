package config

import "net"

// BuildConfiguration represents parameters for building a VM image
type BuildConfiguration struct {
	KernelImage   string
	KernelCmdLine string
	Metadata      string
	CopyFiles     []string

	ID   string `json:"-"` // ID is automatically generated
	Name string `json:"name"`
	// Source can either be a Docker image, a tar file, a folder, or in the future a file or partition containing an ext4 filesystem
	Source    string    `json:"source"`
	Resources Resources `json:"resources"`
}

type Resources struct {
	CPU    CPUResources    `json:"cpu"`
	Memory MemoryResources `json:"memory"`
	Drives DriveResources  `json:"drives"`
}

type CPUResources struct {
	VCPUs          uint32 `json:"vCPUs"` // Count?
	Template       string `json:"template"`
	Hyperthreading bool   `json:"hyperthreading"`
}

type MemoryResources struct {
	RAM uint32 `json:"ram"` // RAMMB, Size or MBSize? TODO: Make this a string and support stuff like 128 M(B), or 8 G(B)
}

type DriveResources struct {
	Root  Drive   `json:"root"`
	Extra []Drive `json:"extra"`
}

type NetworkResources struct {
	// TAP is a slice of tap adapters to create and connect to the VM
	TAP []TAPAdapter `json:"tap"`
	// Vsock
	Vsock []Vsock `json:"vsock"`
}

type TAPAdapter struct {
	// Name of the tap adapter on the host
	Name string `json:"name"`
	// MAC is the MAC address of the tap interface inside of the VM
	MAC net.HardwareAddr `json:"mac"`
	// HostBridge specifies which bridge to connect the tap adapter to
	HostBridge string `json:"hostBridge"`
}

type Vsock struct {
	// Path specifies the path to the socket on the host
	Path string `json:"path"`
	// CID is the client ID inside of the VM
	CID uint32 `json:"mac"`
}

type Drive struct {
	// Path can either be something like /dev/sda1, a file with an ext4 filesystem in it, or a label
	Path     string  `json:"path"`
	ReadOnly *bool   `json:"readOnly"`
	UUID     *string `json:"uuid"`
}
