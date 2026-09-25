package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tidykit/scanner"
)

type DynamicMetric struct {
	Title       string
	Bytes       uint64
	IsLoading   bool
	LabelWidget *widget.Label
	BtnWidget   *widget.Button
	Items       []scanner.CleanableItem
	CleanFunc   func() error
}

func MakeDashboardView(win fyne.Window) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Dashboard", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	metrics := []*DynamicMetric{
		{Title: "Purgable Cache"},
		{Title: "Large Files"},
		{Title: "Duplicates"},
		{Title: "Trash"},
	}

	cardsGrid := container.NewGridWithColumns(4)

	for _, m := range metrics {
		m.IsLoading = true
		m.LabelWidget = widget.NewLabel("Scanning...")

		metricRef := m
		m.BtnWidget = widget.NewButton("Check & Optimize", func() {
			if len(metricRef.Items) > 0 {
				ShowCleanModal(metricRef.Title, metricRef.Items, metricRef.CleanFunc, win.Canvas())
			}
		})
		m.BtnWidget.Disable()

		cardContent := container.NewVBox(
			widget.NewLabelWithStyle(m.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			m.LabelWidget,
			m.BtnWidget,
		)
		cardsGrid.Add(widget.NewCard("", "", cardContent))
	}

	actionSectionTitle := widget.NewLabelWithStyle("Quick Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var trashMetric *DynamicMetric
	for _, m := range metrics {
		if m.Title == "Trash" {
			trashMetric = m
			break
		}
	}

	actionsGrid := container.NewGridWithColumns(3,
		widget.NewButtonWithIcon("Deep Scan System", theme.SearchIcon(), func() {
			ShowScanModal("/", win.Canvas())
		}),
		widget.NewButtonWithIcon("Empty Trash", theme.DeleteIcon(), func() {
			if trashMetric != nil {
				_ = scanner.EmptyTrash()
				items, _ := scanner.ScanTrashFiles()
				fyne.Do(func() {
					total := uint64(0)
					for _, item := range items {
						total += uint64(item.Size)
					}
					trashMetric.Bytes = total
					trashMetric.LabelWidget.SetText(scanner.FormatBytes(total))
				})
			}
		}),
		widget.NewButtonWithIcon("Analyze Disk Usage", theme.StorageIcon(), func() {
			fmt.Println("Triggered: Disk Analysis")
		}),
		widget.NewButtonWithIcon("Manage Startup Apps", theme.ComputerIcon(), func() {
			fmt.Println("Triggered: Startup Apps")
		}),
		widget.NewButtonWithIcon("View System Logs", theme.DocumentIcon(), func() {
			fmt.Println("Triggered: View Logs")
		}),
		widget.NewButtonWithIcon("Optimize /tmp Dir", theme.DocumentIcon(), func() {
			fmt.Println("Triggered: Optimize /tmp")
		}),
	)

	topContainer := container.NewVBox(
		header,
		widget.NewSeparator(),
	)

	mainBody := container.NewVBox(
		cardsGrid,
		widget.NewSeparator(),
		actionSectionTitle,
		actionsGrid,
	)

	scrollContent := container.NewScroll(mainBody)

	go runBackgroundScan(metrics)

	return container.NewBorder(topContainer, nil, nil, nil, scrollContent)
}

func runBackgroundScan(metrics []*DynamicMetric) {
	for _, m := range metrics {
		go func(metric *DynamicMetric) {
			var size uint64

			switch metric.Title {
			case "Purgable Cache":
				items, err := scanner.ScanCacheFiles()
				if err == nil {
					for _, item := range items {
						size += uint64(item.Size)
					}
					metric.Items = items
					metric.CleanFunc = func() error {
						for _, item := range items {
							_ = os.Remove(item.Path)
						}
						return nil
					}
				}
			case "Large Files":
				files, err := scanner.ScanLargeFiles(100*1024*1024, 100)
				if err == nil {
					for _, f := range files {
						size += uint64(f.Size)
					}
					metric.Items = files
					metric.CleanFunc = func() error {
						for _, item := range files {
							_ = os.Remove(item.Path)
						}
						return nil
					}
				}
			case "Duplicates":
				items, err := scanner.ScanDuplicates()
				if err == nil {
					for _, item := range items {
						size += uint64(item.Size)
					}
					metric.Items = items
					metric.CleanFunc = func() error {
						for _, item := range items {
							_ = os.Remove(item.Path)
						}
						return nil
					}
				}
			case "Trash":
				items, err := scanner.ScanTrashFiles()
				if err == nil {
					for _, item := range items {
						size += uint64(item.Size)
					}
					metric.Items = items
					metric.CleanFunc = func() error {
						return scanner.EmptyTrash()
					}
				}
			}

			fyne.Do(func() {
				metric.Bytes = size
				metric.IsLoading = false
				metric.LabelWidget.SetText(scanner.FormatBytes(size))
				metric.BtnWidget.Enable()
			})
		}(m)
	}
}
