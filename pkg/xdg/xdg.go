package xdg

// REFERENCE: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envCacheHome = "XDG_CACHE_HOME"

	envConfigDirs = "XDG_CONFIG_DIRS"
	envConfigHome = "XDG_CONFIG_HOME"

	envDataDirs = "XDG_DATA_DIRS"
	envDataHome = "XDG_DATA_HOME"

	envRuntimeDir = "XDG_RUNTIME_DIR"
)

func envOrDefault(env string, defaultF func() (string, error)) (string, error) {
	if res, ok := os.LookupEnv(env); ok {
		if !filepath.IsAbs(res) {
			return "", fmt.Errorf("path '%s' is not absolute", res)
		}

		return res, nil
	}
	return defaultF()
}

// CacheHome TODO.
func CacheHome() (string, error) { return envOrDefault(envCacheHome, defaultCacheHome) }

// ConfigDirs TODO.
func ConfigDirs() ([]string, error) {
	toSplit, err := envOrDefault(envConfigDirs, defaultConfigDirs)
	if err != nil {
		return nil, err
	}

	return strings.Split(toSplit, ":"), nil
}

// ConfigHome TODO.
func ConfigHome() (string, error) { return envOrDefault(envConfigHome, defaultConfigHome) }

// ConfigSearchDirs TODO.
func ConfigSearchDirs() ([]string, error) {
	home, err := ConfigHome()
	if err != nil {
		return nil, err
	}

	dirs, err := ConfigDirs()
	if err != nil {
		return nil, err
	}

	return append([]string{home}, dirs...), nil
}

// DataDirs TODO.
func DataDirs() ([]string, error) {
	toSplit, err := envOrDefault(envDataDirs, defaultDataDirs)
	if err != nil {
		return nil, err
	}

	return strings.Split(toSplit, ":"), nil
}

// DataHome TODO.
func DataHome() (string, error) { return envOrDefault(envDataHome, defaultDataHome) }

// DataSearchDirs TODO.
func DataSearchDirs() ([]string, error) {
	home, err := DataHome()
	if err != nil {
		return nil, err
	}

	dirs, err := DataDirs()
	if err != nil {
		return nil, err
	}

	return append([]string{home}, dirs...), nil
}

// RuntimeDir TODO.
func RuntimeDir() (string, error) { return envOrDefault(envRuntimeDir, defaultRuntimeDir) }
