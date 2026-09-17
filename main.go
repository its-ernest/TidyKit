package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"tidykit/ui"
)

func main() {
	a := app.New()
	w := a.NewWindow("TidyKit")
	w.Resize(fyne.NewSize(900, 560))

	w.SetContent(ui.BuildMainUI(w))
	w.ShowAndRun()
}
