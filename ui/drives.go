package ui

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tidykit/scanner"
)

type DriveInfo = scanner.DriveInfo

func FetchDrives() ([]DriveInfo, error) {
	return scanner.FetchDrives()
}

func MakeDrivesView(win fyne.Window) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Connected Drives", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	drives, err := FetchDrives()
	if err != nil {
		return container.NewVBox(
			header,
			widget.NewSeparator(),
			widget.NewLabel(fmt.Sprintf("Failed to inspect drives: %v", err)),
		)
	}

	topContainer := container.NewVBox(
		header,
		widget.NewSeparator(),
	)

	cardsGrid := container.NewGridWithColumns(2)
	for _, drive := range drives {
		cardsGrid.Add(createDriveCard(drive, win))
	}

	scrollableGrid := container.NewScroll(container.NewVBox(cardsGrid))

	return container.NewBorder(topContainer, nil, nil, nil, scrollableGrid)
}

func createDriveCard(drive DriveInfo, win fyne.Window) fyne.CanvasObject {
	usedGB := float64(drive.UsedBytes) / (1024 * 1024 * 1024)
	totalGB := float64(drive.TotalBytes) / (1024 * 1024 * 1024)
	usedRatio := float64(drive.UsedBytes) / float64(drive.TotalBytes)

	title := fmt.Sprintf("%s (%s)", drive.MountPoint, drive.Device)
	capacityText := fmt.Sprintf("%.1f GB / %.1f GB", usedGB, totalGB)

	prog := widget.NewProgressBar()
	prog.SetValue(usedRatio)

	// Action 1: Open drive in native file manager (Linux default: xdg-open)
	openBtn := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		_ = exec.Command("xdg-open", drive.MountPoint).Start()
	})

	// Action 2: Scan drive and show modal file tree
	scanBtn := widget.NewButtonWithIcon("Scan", theme.SearchIcon(), func() {
		ShowScanModal(drive.MountPoint, win.Canvas())
	})

	buttonRow := container.NewGridWithColumns(2, openBtn, scanBtn)

	cardContent := container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(capacityText),
		prog,
		buttonRow,
	)

	return widget.NewCard("", "", cardContent)
}
