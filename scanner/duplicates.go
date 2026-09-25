package scanner

import (
	"crypto/md5"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type DuplicateGroup struct {
	Files []string
	Size  int64
}

func ScanDuplicates() ([]CleanableItem, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sizeMap := make(map[int64][]string)

	_ = filepath.WalkDir(home, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".cache" || name == ".local" || name == ".npm" || name == ".cargo" || name == ".gnupg" || name == ".config" && filepath.Base(filepath.Dir(path)) == ".config" {
				return fs.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.Size() > 0 {
			sizeMap[info.Size()] = append(sizeMap[info.Size()], path)
		}

		return nil
	})

	hashMap := make(map[string]*DuplicateGroup)
	var groups []*DuplicateGroup

	for _, files := range sizeMap {
		if len(files) < 2 {
			continue
		}

		for _, filePath := range files {
			h, err := md5File(filePath)
			if err != nil {
				continue
			}

			key := h
			if g, ok := hashMap[key]; ok {
				g.Files = append(g.Files, filePath)
			} else {
				g := &DuplicateGroup{
					Files: []string{filePath},
					Size:  0,
				}
				hashMap[key] = g
				groups = append(groups, g)
			}
		}
	}

	var result []CleanableItem
	for _, g := range groups {
		if len(g.Files) > 1 {
			info, _ := os.Stat(g.Files[0])
			if info != nil {
				g.Size = info.Size() * int64(len(g.Files)-1)
			}
			for i := 1; i < len(g.Files); i++ {
				result = append(result, CleanableItem{
					Name: g.Files[i],
					Path: g.Files[i],
					Size: g.Size / int64(len(g.Files)-1),
				})
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	return result, nil
}

func md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
