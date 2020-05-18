package version

import (
	"fmt"
	"runtime"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	igniteruntime "github.com/weaveworks/ignite/pkg/runtime"
)

// These variables are used for building ignite with ldflag overrides
var (
	gitMajor           = ""
	gitMinor           = ""
	gitVersion         = ""
	gitCommit          = ""
	gitTreeState       = ""
	buildDate          = ""
	firecrackerVersion = ""

	// allow overriding the DEFAULT_SANDBOX_IMAGE_* constants
	sandboxImageName      = ""
	sandboxImageTag       = ""
	sandboxImageDelimeter = ":" // set to "@" to support using sha256:<hash> as a "tag"

	// allow overriding the DEFAULT_KERNEL_IMAGE_* constants
	kernelImageName      = ""
	kernelImageTag       = ""
	kernelImageDelimeter = ":" // set to "@" to support using sha256:<hash> as a "tag"
)

// Image represents an OCI image
// TODO: use a shared or upstream OCI Image type
type Image struct {
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Delimeter string `json:"delimeter"`
}

// String returns Image{} as a human-friendly version string.
func (image Image) String() string {
	return fmt.Sprintf("%s%s%s", image.Name, image.Delimeter, image.Tag)
}

// Info stores information about a component's version
type Info struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
	SandboxImage Image  `json:"sandboxImage"`
	KernelImage  Image  `json:"kernelImage"`
}

// String returns Info{} as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}

// GetIgnite gets ignite's version
func GetIgnite() Info {
	info := Info{
		Major:        gitMajor,
		Minor:        gitMinor,
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),

		SandboxImage: Image{
			Name:      constants.DEFAULT_SANDBOX_IMAGE_NAME,
			Tag:       constants.DEFAULT_SANDBOX_IMAGE_TAG,
			Delimeter: sandboxImageDelimeter,
		},

		KernelImage: Image{
			Name:      constants.DEFAULT_KERNEL_IMAGE_NAME,
			Tag:       constants.DEFAULT_KERNEL_IMAGE_TAG,
			Delimeter: kernelImageDelimeter,
		},
	}

	if sandboxImageName != "" {
		info.SandboxImage.Name = sandboxImageName
	}
	if sandboxImageTag != "" {
		info.SandboxImage.Tag = sandboxImageTag
	}

	if kernelImageName != "" {
		info.KernelImage.Name = kernelImageName
	}
	if kernelImageTag != "" {
		info.KernelImage.Tag = kernelImageTag
	}

	return info
}

// GetFirecracker returns firecracker's version
func GetFirecracker() Info {
	return Info{
		GitVersion: firecrackerVersion,
	}
}

// GetCurrentRuntime returns the current configured runtime
func GetCurrentRuntime() igniteruntime.Name {
	return providers.RuntimeName
}
