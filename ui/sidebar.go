package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MakeSidebar builds the navigation panel on the left.
func MakeSidebar(onSelect func(viewID string)) fyne.CanvasObject {
	dashBtn := widget.NewButton("Dashboard", func() {
		onSelect("dashboard")
	})

	drivesBtn := widget.NewButton("Connected Drives", func() {
		onSelect("drives")
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Menu", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		dashBtn,
		drivesBtn,
	)
}
