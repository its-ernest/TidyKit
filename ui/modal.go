package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tidykit/scanner"
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

	sizedContainer := container.NewPadded(modalLayout)

	modal = widget.NewModalPopUp(sizedContainer, canvas)
	modal.Resize(fyne.NewSize(750, 500))
	modal.Show()

	return modal
}

type TreeItem = scanner.TreeItem

// Helper to format byte sizes into human-readable strings
func formatBytes(bytes uint64) string {
	return scanner.FormatBytes(bytes)
}

// ShowScanModal launches a recursive directory scanner and displays results in a navigable list.
func ShowScanModal(rootPath string, canvas fyne.Canvas) {
	statusLabel := widget.NewLabel(fmt.Sprintf("Scanning %s...", rootPath))
	progress := widget.NewProgressBarInfinite()

	loadingView := container.NewCenter(
		container.NewVBox(
			statusLabel,
			progress,
		),
	)

	contentBox := container.NewStack(loadingView)
	_ = ShowModal(fmt.Sprintf("Scan Results: %s", rootPath), contentBox, canvas)

	go func() {
		nodes, err := scanner.ScanDir(rootPath)
		if err != nil {
			fyne.Do(func() {
				contentBox.Objects = []fyne.CanvasObject{
					widget.NewLabel("Failed to access or read directory."),
				}
				contentBox.Refresh()
			})
			return
		}

		cleanRoot := filepath.Clean(rootPath)
		var currentPath string = cleanRoot

		pathLabel := widget.NewLabel(currentPath)

		var list *widget.List
		list = widget.NewList(
			func() int {
				if item, ok := nodes[currentPath]; ok {
					return len(item.Children)
				}
				return 0
			},
			func() fyne.CanvasObject {
				deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
				return container.NewBorder(
					nil, nil,
					container.NewHBox(
						widget.NewIcon(theme.FolderIcon()),
						widget.NewLabel("Template"),
					),
					deleteBtn,
				)
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				border := obj.(*fyne.Container)
				leftContent := border.Objects[0].(*fyne.Container)
				deleteBtn := border.Objects[1].(*widget.Button)

				icon := leftContent.Objects[0].(*widget.Icon)
				label := leftContent.Objects[1].(*widget.Label)

				if item, ok := nodes[currentPath]; ok {
					if id < len(item.Children) {
						childPath := item.Children[id]
						if childItem, ok := nodes[childPath]; ok {
							if childItem.IsDir {
								icon.SetResource(theme.FolderIcon())
								label.SetText(fmt.Sprintf("%s/", childItem.Name))
							} else {
								icon.SetResource(theme.FileIcon())
								label.SetText(fmt.Sprintf("%s (%s)", childItem.Name, scanner.FormatBytes(uint64(childItem.Size))))
							}

							childPathCopy := childPath
							deleteBtn.OnTapped = func() {
								go func() {
									_ = os.RemoveAll(childPathCopy)
									fyne.Do(func() {
										removeNodeAndDescendants(nodes, childPathCopy)
										if parentItem, ok := nodes[currentPath]; ok {
											for i, child := range parentItem.Children {
												if child == childPathCopy {
													parentItem.Children = append(parentItem.Children[:i], parentItem.Children[i+1:]...)
													break
												}
											}
										}
										list.Refresh()
									})
								}()
							}
						}
					}
				}
			},
		)

		backBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			parent := filepath.Dir(currentPath)
			if parent != currentPath {
				currentPath = parent
				pathLabel.SetText(currentPath)
				list.Refresh()
			}
		})
		if currentPath == cleanRoot {
			backBtn.Disable()
		}

		list.OnSelected = func(id widget.ListItemID) {
			if item, ok := nodes[currentPath]; ok {
				if id < len(item.Children) {
					childPath := item.Children[id]
					if childItem, ok := nodes[childPath]; ok && childItem.IsDir {
						currentPath = childPath
						pathLabel.SetText(currentPath)
						if filepath.Dir(currentPath) != cleanRoot {
							backBtn.Enable()
						} else {
							backBtn.Disable()
						}
						list.Refresh()
					}
				}
			}
		}

		header := container.NewBorder(
			nil, nil,
			widget.NewLabelWithStyle("Current Directory", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			backBtn,
			pathLabel,
		)

		content := container.NewBorder(
			container.NewVBox(
				header,
				widget.NewSeparator(),
			),
			nil, nil, nil,
			list,
		)

		fyne.Do(func() {
			contentBox.Objects = []fyne.CanvasObject{content}
			contentBox.Refresh()
		})
	}()
}

func removeNodeAndDescendants(nodes map[string]*scanner.TreeItem, path string) {
	if item, ok := nodes[path]; ok && item.IsDir {
		for _, child := range item.Children {
			removeNodeAndDescendants(nodes, child)
		}
	}
	delete(nodes, path)
}
