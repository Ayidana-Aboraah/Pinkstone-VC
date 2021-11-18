package main

import (
	"business.go/Data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var (
	mainMenu = fyne.NewContainer()
	dataMenu = fyne.NewContainer()
	shopMenu = fyne.NewContainer()
	itemMenu = fyne.NewContainer()
)

func main() {
	a := app.NewWithID("Bronze Hermes")
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
		widget.NewButton("Debug", func() {
			a.SendNotification(fyne.NewNotification(Data.ConvertDate(time.Now()), Data.ConvertClock(time.Now())))
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

	itemMenu = container.NewVBox(
		widget.NewLabelWithStyle("", fyne.TextAlign(1), fyne.TextStyle{Italic: true}),
		//widget.NewLabel(strconv.Itoa(tempSale.id)),
		widget.NewEntry(),//Name
		widget.NewEntry(),//Price
		widget.NewEntry(),//Cost
		widget.NewEntry(),//Inventory
		widget.NewButton("Submit", func() {
			//Do what ever submitting data does
		}),
		widget.NewButton("Cancel", func() {
			//tempSale = nil
			w.SetContent(mainMenu)
		}),
	)

	w.SetContent(mainMenu)

	w.ShowAndRun()
}