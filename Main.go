package main

import (
	"business.go/Cam"
	"business.go/Data"
	"fmt"
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
			widget.NewLabel("Shopping"),
			widget.NewButton("Time", func() {
				a.SendNotification(fyne.NewNotification(Data.ConvertDate(time.Now()), Data.ConvertClock(time.Now())))
			}),
			widget.NewButton("", func() {

			}),
		)),
		
			container.NewTabItem("Barcodes", container.NewVBox(
				testTitle,
				//widget.NewCard("Homies", "You Thought...", widget.NewEntry()),
				widget.NewButton("Camera", func(){
					//Cam.OpenCam()
				}),
				widget.NewButton("Barcode 01", func (){
					file, _ := os.Open(Cam.Path + "Online Test 01.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 02", func (){
					file, _ := os.Open(Cam.Path + "Online Test 02.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 03", func (){
					file, _ := os.Open(Cam.Path + "Online Test 03.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 04", func (){
					file, _ := os.Open(Cam.Path + "Test01.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 05", func (){
					file, _ := os.Open(Cam.Path + "Online Test 05.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Add Barcode 05 To DataBase", func() {
					file, _ := os.Open(Cam.Path + "Online Test 05.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).GetText()

					fmt.Println(id)
					//Add this to the pop up menu when I can do that
					//testTitle.SetText("ID: " + id + " added to Test Data Base")
				}),
			)),
			
			container.NewTabItem("Info", container.NewVBox(

				container.NewHSplit(
					widget.NewLabel("Test 3"),
					container.NewVScroll(
						widget.NewButton("Test 01", func() {
							//
						}),
					)),

					widget.NewForm(
						widget.NewFormItem("Id", widget.NewLabel("ID")),
						widget.NewFormItem("Price", widget.NewEntry()),
						widget.NewFormItem("Cost", widget.NewEntry()),
						widget.NewFormItem("Inventory", widget.NewEntry()),
					),
			)),

		),
	)

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func CreateNewItem(){

}