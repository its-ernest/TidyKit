package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ScanTrashFiles() ([]CleanableItem, error) {
	trashDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "Trash", "files")
	var items []CleanableItem
	err := filepath.WalkDir(trashDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err == nil && !d.IsDir() {
			items = append(items, CleanableItem{
				Name: filepath.Base(path),
				Path: path,
				Size: info.Size(),
			})
		}
		return nil
	})
	return items, err
}

func EmptyTrash() error {
	trashDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "Trash")
	_ = os.RemoveAll(trashDir)
	return os.MkdirAll(trashDir, 0755)
}
