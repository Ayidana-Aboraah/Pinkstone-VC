package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Bronze Hermes")
	//a.SetIcon(icon)

	title := widget.NewLabel("Welcome!")
	title.Alignment = fyne.TextAlign(1)
	w.SetContent(container.NewVBox(
		title,
		widget.NewButton("Data", func() {
			title.SetText("Welcome :)")
			//Change To Data Menu
		}),
		widget.NewButton("Camera", func(){
			//StartCam()
		}),
		widget.NewButton("Quit", func(){
			w.Close()
		}),
	))
	w.ShowAndRun()
}