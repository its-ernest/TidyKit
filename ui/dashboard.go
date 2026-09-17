package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
		m.BtnWidget = widget.NewButton("Clean Now", func() {
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
	// 1. Perform background processing / heavy file checks here
	time.Sleep(1200 * time.Millisecond)

	mockSizes := []uint64{
		2 * 1024 * 1024 * 1024,  // 2.0 GB
		14 * 1024 * 1024 * 1024, // 14.0 GB
		850 * 1024 * 1024,       // 850 MB
		120 * 1024 * 1024,       // 120 MB
	}

	for i, m := range metrics {
		m.Bytes = mockSizes[i]
		m.IsLoading = false
		formattedSize := formatBytes(m.Bytes)

		// 2. Capture local variables for closure safely
		label := m.LabelWidget
		btn := m.BtnWidget

		// 3. Dispatch UI mutations onto Fyne's main thread using fyne.Do
		fyne.Do(func() {
			label.SetText(formattedSize)
			btn.Enable()
		})
	}
}
