package config

import (
	"os"
	"path/filepath"

	"github.com/oligarch316/go-skeleton/pkg/xdg"
)

/* NOTE:
   This is all MVP. Notable shorcomings include...
   - Not really to XDG spec given we ignore XDG_CONFIG_DIRS
   - type XDGGlobs support ?
   - Flexibility between OneOf vs AllOf vs etc. behavior ?
*/

// XDGPaths TODO.
type XDGPaths struct{ RelativePaths []FilePath }

// Readers TODO.
func (xp XDGPaths) Readers() ([]Reader, error) {
	searchDirs, err := xdg.ConfigSearchDirs()
	if err != nil {
		return nil, err
	}

	var pathSources []Source
	for _, searchDir := range searchDirs {
		for _, relPathSource := range xp.RelativePaths {
			fullPath := filepath.Join(searchDir, relPathSource.Path)

			info, err := os.Stat(fullPath)
			switch {
			case os.IsNotExist(err):
				continue
			case err != nil:
				return nil, err
			case info.IsDir():
				continue
			}

			pathSources = append(pathSources, FilePath{
				Format: relPathSource.Format,
				Path:   fullPath,
			})
		}
	}

	return OneOf(pathSources).Readers()
}
