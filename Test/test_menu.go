package test

import (
	"BronzeHermes/Database"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func TestMenu(a fyne.App, w fyne.Window) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("Display Database", func() {
			dialog.ShowInformation("Databases", fmt.Sprint(Database.Databases), w)
			dialog.ShowInformation("Name Keys", fmt.Sprint(Database.NameKeys), w)
		}))
}
