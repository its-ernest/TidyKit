package ui

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ShowModal wraps arbitrary CanvasObject content inside a dismissible modal overlay.
func ShowModal(title string, content fyne.CanvasObject, canvas fyne.Canvas) *widget.PopUp {
	var modal *widget.PopUp

	header := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	closeBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		if modal != nil {
			modal.Hide()
		}
	})

	topBar := container.NewBorder(nil, nil, header, closeBtn)
	modalLayout := container.NewBorder(
		container.NewVBox(topBar, widget.NewSeparator()),
		nil, nil, nil,
		content,
	)

	// Fixed size overlay container
	padded := container.NewPadded(modalLayout)
	sizedContainer := container.NewGridWithRows(1, padded)

	modal = widget.NewModalPopUp(sizedContainer, canvas)
	modal.Resize(fyne.NewSize(700, 500))
	modal.Show()

	return modal
}

// TreeItem represents a file or folder in our recursive scan model
type TreeItem struct {
	Path     string
	Name     string
	Size     int64
	IsDir    bool
	Children []string // Child paths
}

// ShowScanModal launches a recursive directory scanner and displays results in a collapsible tree.
func ShowScanModal(rootPath string, canvas fyne.Canvas) {
	statusLabel := widget.NewLabel("Scanning file tree...")
	progress := widget.NewProgressBarInfinite()

	loadingView := container.NewVBox(
		statusLabel,
		progress,
	)

	contentBox := container.NewStack(loadingView)
	modal := ShowModal(fmt.Sprintf("Scan Results: %s", rootPath), contentBox, canvas)

	go func() {
		nodes := make(map[string]*TreeItem)

		// Recursive scan capped at reasonable depth for quick UX display
		_ = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip unreadable paths
			}

			info, err := d.Info()
			var size int64
			if err == nil {
				size = info.Size()
			}

			item := &TreeItem{
				Path:  path,
				Name:  d.Name(),
				Size:  size,
				IsDir: d.IsDir(),
			}

			nodes[path] = item

			parent := filepath.Dir(path)
			if parent != path {
				if parentItem, ok := nodes[parent]; ok {
					parentItem.Children = append(parentItem.Children, path)
				}
			}

			return nil
		})

		// Build Collapsible Recursive Tree Widget
		tree := widget.NewTree(
			func(id string) []string {
				if id == "" {
					return []string{rootPath}
				}
				if item, ok := nodes[id]; ok {
					return item.Children
				}
				return nil
			},
			func(id string) bool {
				if item, ok := nodes[id]; ok {
					return item.IsDir
				}
				return false
			},
			func(branch bool) fyne.CanvasObject {
				return container.NewHBox(
					widget.NewIcon(theme.FileIcon()),
					widget.NewLabel("Template Path Item"),
				)
			},
			func(id string, branch bool, obj fyne.CanvasObject) {
				box := obj.(*fyne.Container)
				icon := box.Objects[0].(*widget.Icon)
				label := box.Objects[1].(*widget.Label)

				item, ok := nodes[id]
				if !ok {
					return
				}

				if branch {
					icon.SetResource(theme.FolderIcon())
					label.SetText(item.Name)
				} else {
					icon.SetResource(theme.FileIcon())
					label.SetText(fmt.Sprintf("%s (%s)", item.Name, formatBytes(uint64(item.Size))))
				}
			},
		)

		fyne.Do(func() {
			contentBox.Objects = []fyne.CanvasObject{tree}
			contentBox.Refresh()
			_ = modal
		})
	}()
}
