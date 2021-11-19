package main

import (
	"business.go/Cam"
	"business.go/Data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/png"
	"os"
	"time"
)

var (
	mainMenu = fyne.NewContainer()
	dataMenu = fyne.NewContainer()
	shopMenu = fyne.NewContainer()
	itemMenu = fyne.NewContainer()

	testMenu = fyne.NewContainer()
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
		widget.NewButton("Test", func() {
			w.SetContent(testMenu)
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
		widget.NewEntry(), //Name
		widget.NewEntry(), //Price
		widget.NewEntry(), //Cost
		widget.NewEntry(), //Inventory
		widget.NewButton("Submit", func() {
			//Do what ever submitting data does
		}),
		widget.NewButton("Cancel", func() {
			//tempSale = nil
			w.SetContent(mainMenu)
		}),
	)

	testTitle := widget.NewLabel("Test 2")
	testMenu = container.NewVBox(
		container.NewAppTabs(container.NewTabItem("Shop", container.NewVBox(
			widget.NewLabel("Test"),
			widget.NewButton("Blue", func() {
				a.SendNotification(fyne.NewNotification(Data.ConvertDate(time.Now()), Data.ConvertClock(time.Now())))
			}),
		)),
			container.NewTabItem("Camera", container.NewVBox(
				testTitle,
				widget.NewButton("Test Image 1", func (){
					file, _ := os.Open(Cam.Path + "Online Test.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText(id)
				}),
				widget.NewButton("Test Image 2", func (){
					file, _ := os.Open(Cam.Path + "test 1.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText(id)
				}),
			)),
			container.NewTabItem("Info", container.NewVBox(
				widget.NewLabel("Test 3"),
			)),
		),
	)

	w.SetContent(mainMenu)
	w.ShowAndRun()
}