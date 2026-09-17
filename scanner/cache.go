package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ScanCacheSize() (int64, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}

	cachePaths := []string{
		filepath.Join(home, ".cache"),
		"/tmp",
	}

	seen := make(map[string]bool)
	var total int64

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
				total += info.Size()
			}
			return nil
		})
	}

	return total, nil
}
