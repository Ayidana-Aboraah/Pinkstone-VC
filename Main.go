package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	mainMenu = fyne.NewContainer()
	dataMenu = fyne.NewContainer()
	shopMenu = fyne.NewContainer()
)

func main() {
	a := app.New()
	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")
	//a.SetIcon(icon)

	title := widget.NewLabel("Welcome!")
	title.Alignment = fyne.TextAlign(1)

	mainMenu = container.NewVBox(
		title,
		widget.NewButton("Checkout", func() {
			w.SetContent(shopMenu)
		}),
		widget.NewButton("Data", func() {
			title.SetText("Data")
			//Change To Data Menu
			w.SetContent(dataMenu)
		}),
		widget.NewButton("Camera", func() {
			//OpenCam()
		}),
		widget.NewButton("Quit", func() {
			w.Close()
		}),
	)

	dataMenu = container.NewVBox(
		title,
		widget.NewButton("Back", func() {
			title.SetText("Bronze Hermes")
			w.SetContent(mainMenu)
		}),
	)

	w.SetContent(mainMenu)

	w.ShowAndRun()
}