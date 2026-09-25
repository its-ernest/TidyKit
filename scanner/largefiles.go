package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type LargeFileInfo struct {
	Path string
	Size int64
}

func ScanLargeFiles(minSize int64, maxResults int) ([]CleanableItem, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	scanRoots := []string{
		home,
	}

	if stat, err := os.Stat("/"); err == nil && stat.IsDir() {
		scanRoots = append(scanRoots, "/")
	}

	var files []CleanableItem
	seen := make(map[string]bool)

	for _, root := range scanRoots {
		if seen[root] {
			continue
		}
		seen[root] = true

		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				name := d.Name()
				if name == "proc" || name == "sys" || name == "dev" || name == "run" || name == "snap" {
					return fs.SkipDir
				}
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			if info.Size() >= minSize {
				files = append(files, CleanableItem{
					Name: path,
					Path: path,
					Size: info.Size(),
				})
			}

			return nil
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > maxResults {
		files = files[:maxResults]
	}

	return files, nil
}
