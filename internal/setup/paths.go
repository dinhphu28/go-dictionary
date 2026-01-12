package setup

import (
	"os"
	"path/filepath"
)

type Paths struct {
	BinDir    string
	ConfigDir string
	DataDir   string
}

func DefaultPaths() Paths {
	home, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()

	return Paths{
		BinDir:    filepath.Join(home, ".local", "bin"),
		ConfigDir: filepath.Join(configDir, "dictionary"),
		DataDir:   filepath.Join(home, ".local", "share", "dictionary"),
	}
}
