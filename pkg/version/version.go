package version

import (
	"fmt"
	"runtime"
)

var (
	gitMajor = ""
	gitMinor = ""
	gitVersion = ""
	gitCommit = ""
	gitTreeState = ""
	buildDate = ""
)

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
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}

// GetIgnite gets ignite's version
func GetIgnite() Info {
	return Info{
		Major: gitMajor,
		Minor: gitMinor,
		GitVersion: gitVersion,
		GitCommit: gitCommit,
		GitTreeState: gitTreeState,
		BuildDate: buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetFirecracker returns firecracker's version
func GetFirecracker() Info {
	return Info{
		Major: "0",
		Minor: "16",
		GitVersion: "v0.16.0",
	}
}
