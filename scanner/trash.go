package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ScanTrashSize() (int64, error) {
	trashDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "Trash", "files")
	var total int64
	err := filepath.WalkDir(trashDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err == nil && !d.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
}
