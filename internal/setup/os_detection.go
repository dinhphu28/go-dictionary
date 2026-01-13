package setup

import "runtime"

type OSInfo struct {
	Name           string
	IsLinux        bool
	IsMac          bool
	IsWindows      bool
	SupportsNative bool
}

func DetectOS() OSInfo {
	os := runtime.GOOS

	return OSInfo{
		Name:           os,
		IsLinux:        os == "linux",
		IsMac:          os == "darwin",
		IsWindows:      os == "windows",
		SupportsNative: os == "linux" || os == "darwin" || os == "windows",
	}
}
