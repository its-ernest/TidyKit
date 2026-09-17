package ui

import (
	"fmt"

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
}

func MakeDashboardView() fyne.CanvasObject {
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
			fmt.Printf("Cleaning %s...\n", metricRef.Title)
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
	actionsGrid := container.NewGridWithColumns(3,
		widget.NewButtonWithIcon("Deep Scan System", theme.SearchIcon(), func() {
			fmt.Println("Triggered: Deep Scan")
		}),
		widget.NewButtonWithIcon("Empty Trash", theme.DeleteIcon(), func() {
			fmt.Println("Triggered: Empty Trash")
		}),
		widget.NewButtonWithIcon("Optimize Storage", theme.SettingsIcon(), func() {
			fmt.Println("Triggered: Optimize Storage")
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

	// Run heavy disk scanning off the UI thread
	go runBackgroundScan(metrics)

	return container.NewBorder(topContainer, nil, nil, nil, scrollContent)
}

func runBackgroundScan(metrics []*DynamicMetric) {
	for _, m := range metrics {
		go func(metric *DynamicMetric) {
			var size uint64

			switch metric.Title {
			case "Purgable Cache":
				cacheSize, err := scanner.ScanCacheSize()
				if err == nil {
					size = uint64(cacheSize)
				}
			case "Large Files":
				files, err := scanner.ScanLargeFiles(100*1024*1024, 100)
				if err == nil {
					for _, f := range files {
						size += uint64(f.Size)
					}
				}
			case "Duplicates":
				groups, err := scanner.ScanDuplicates()
				if err == nil {
					for _, g := range groups {
						size += uint64(g.Size)
					}
				}
			case "Trash":
				trashSize, err := scanner.ScanTrashSize()
				if err == nil {
					size = uint64(trashSize)
				}
			}

			fyne.Do(func() {
				metric.Bytes = size
				metric.IsLoading = false
				metric.LabelWidget.SetText(formatBytes(size))
				metric.BtnWidget.Enable()
			})
		}(m)
	}
}
