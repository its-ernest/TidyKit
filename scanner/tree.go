package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
)

type TreeItem struct {
	Path     string
	Name     string
	Size     int64
	IsDir    bool
	Children []string
}

type CleanableItem struct {
	Name string
	Path string
	Size int64
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func ScanDir(rootPath string) (map[string]*TreeItem, error) {
	nodes := make(map[string]*TreeItem)
	cleanRoot := filepath.Clean(rootPath)

	err := filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		info, err := d.Info()
		var size int64
		if err == nil {
			size = info.Size()
		}

		nodes[path] = &TreeItem{
			Path:  path,
			Name:  d.Name(),
			Size:  size,
			IsDir: d.IsDir(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := nodes[cleanRoot]; !ok {
		return nil, fmt.Errorf("failed to access or read directory")
	}

	for path := range nodes {
		if path == cleanRoot {
			continue
		}
		parent := filepath.Dir(path)
		if parentItem, ok := nodes[parent]; ok {
			parentItem.Children = append(parentItem.Children, path)
		}
	}

	for _, item := range nodes {
		if item.IsDir {
			item.Size = calculateDirSize(nodes, item)
		}
	}

	for _, item := range nodes {
		if len(item.Children) > 0 {
			sort.SliceStable(item.Children, func(i, j int) bool {
				childI := nodes[item.Children[i]]
				childJ := nodes[item.Children[j]]
				if childI.IsDir && !childJ.IsDir {
					return true
				}
				if !childI.IsDir && childJ.IsDir {
					return false
				}
				return childI.Size > childJ.Size
			})
		}
	}

	return nodes, nil
}

func calculateDirSize(nodes map[string]*TreeItem, dir *TreeItem) int64 {
	var total int64
	for _, childPath := range dir.Children {
		if child, ok := nodes[childPath]; ok {
			if child.IsDir {
				total += calculateDirSize(nodes, child)
			} else {
				total += child.Size
			}
		}
	}
	return total
}
