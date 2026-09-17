package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Router struct {
	views      map[string]fyne.CanvasObject
	contentBox *fyne.Container
}

func BuildMainUI(win fyne.Window) fyne.CanvasObject {
	router := &Router{
		views:      make(map[string]fyne.CanvasObject),
		contentBox: container.NewStack(),
	}

	router.views["dashboard"] = MakeDashboardView()
	router.views["drives"] = MakeDrivesView(win)

	sidebar := MakeSidebar(func(viewID string) {
		router.NavigateTo(viewID)
	})

	router.NavigateTo("dashboard")

	return container.NewBorder(nil, nil, sidebar, nil, router.contentBox)
}

func (r *Router) NavigateTo(viewID string) {
	if view, exists := r.views[viewID]; exists {
		r.contentBox.Objects = []fyne.CanvasObject{view}
		r.contentBox.Refresh()
	}
}
