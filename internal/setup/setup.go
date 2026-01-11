package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Mode          string `json:"mode"`
	ResourcesPath string `json:"resources_path"`
}

func Run() error {
	cfgPath, err := configPath()
	if err != nil {
		return err
	}

	resPath, err := resourcesPath()
	if err != nil {
		return err
	}

	// Ensure dirs exist
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(resPath, 0755); err != nil {
		return err
	}

	cfg := Config{
		Mode:          "native",
		ResourcesPath: resPath,
	}

	if err := writeConfig(cfgPath, cfg); err != nil {
		return err
	}

	if err := installNativeManifest(); err != nil {
		return err
	}

	fmt.Println("✔ Dictionary setup complete")
	fmt.Println("Config:", cfgPath)
	fmt.Println("Resources:", resPath)
	fmt.Println("Run: dictionary")

	return nil
}

func writeConfig(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "dictionary", "config.json"), nil
}

func resourcesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "dictionary", "resources"), nil
}
