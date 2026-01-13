package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target, info.Mode())
	})
}

func Install(paths Paths) error {
	// 1. Install binary
	binPath := paths.BinPath
	if err := installBinary(binPath); err != nil {
		return err
	}

	if err := installConfigs(paths); err != nil {
		return err
	}

	// 4. Install resources
	resourcesPath := filepath.Join(paths.DataDir, "resources")
	if err := installResources(resourcesPath); err != nil {
		return err
	}

	// 5. Install browser manifests
	installNativeMessagingManifests(binPath)

	return nil
}

func installBinary(path string) error {
	if err := copyFile(
		"./dictionary",
		path,
		0o755,
	); err != nil {
		return fmt.Errorf("install binary: %w", err)
	}
	return nil
}

func installConfigs(paths Paths) error {
	// 1. Install config.json (only if not exists)
	cfgDst := filepath.Join(paths.ConfigDir, "config.json")
	if _, err := os.Stat(cfgDst); os.IsNotExist(err) {
		if err := copyFile("./config.json", cfgDst, 0o644); err != nil {
			return err
		}
	}

	// 2. Install runtime.json
	runtimeCfg := filepath.Join(paths.ConfigDir, "runtime.json")
	if err := os.MkdirAll(paths.ConfigDir, 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(runtimeCfg); os.IsNotExist(err) {
		content := fmt.Sprintf(`{
  "mode": "native",
  "resources_path": "%s"
}
`, filepath.Join(paths.DataDir, "resources"))

		if err := os.WriteFile(runtimeCfg, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func installNativeMessagingManifests(bin string) {
	chrome, firefox := detectBrowsers()

	if chrome {
		manifest := chromeManifest(bin, "kpgiaenkniiaacjbiipbmcdjfbjmgmll")
		installChromeManifest(manifest)
	}

	if firefox {
		manifest := firefoxManifest(bin)
		installFirefoxManifest(manifest)
	}
}

func installResources(resourcesPath string) error {
	if _, err := os.Stat("./resources"); err == nil {
		if err := copyDir("./resources", resourcesPath); err != nil {
			return err
		}
	}
	return nil
}
