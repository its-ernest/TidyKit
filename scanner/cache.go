package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ScanCacheFiles() ([]CleanableItem, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cachePaths := []string{
		filepath.Join(home, ".cache"),
		"/tmp",
	}

	seen := make(map[string]bool)
	var items []CleanableItem

	for _, cacheRoot := range cachePaths {
		if seen[cacheRoot] {
			continue
		}
		seen[cacheRoot] = true

		_ = filepath.WalkDir(cacheRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err == nil {
				items = append(items, CleanableItem{
					Name: path,
					Path: path,
					Size: info.Size(),
				})
			}
			return nil
		})
	}

	return items, nil
}
