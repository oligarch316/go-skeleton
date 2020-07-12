package xdg

// REFERENCE: https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPFileSystem/Articles/WhereToPutFiles.html

import (
	"errors"
	"os"
	"path/filepath"
)

func homeDir() (string, error) {
	if res, ok := os.LookupEnv("HOME"); ok {
		return res, nil
	}
	return "", errors.New("$HOME is not defined")
}

func homeJoin(segs ...string) (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{home}, segs...)...), nil
}

func defaultCacheHome() (string, error) { return homeJoin("Library", "Caches") }

func defaultConfigDirs() (string, error) { return "/etc/xdg", nil }
func defaultConfigHome() (string, error) { return homeJoin("Library", "Application Support") }

func defaultDataDirs() (string, error) { return "/usr/local/share:/usr/share", nil }
func defaultDataHome() (string, error) { return homeJoin("Library", "ApplicationSupport") }

func defaultRuntimeDir() (string, error) { return os.TempDir(), nil }
