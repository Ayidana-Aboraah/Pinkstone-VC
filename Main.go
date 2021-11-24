package main

import (
	"business.go/Cam"
	"business.go/Data"
	"business.go/UI"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/png"
	"os"
	"strconv"
	"time"
)

var ShoppingCart []*Data.Sale

var (
	mainMenu = fyne.NewContainer()
	dataMenu = fyne.NewContainer()
	shopMenu = fyne.NewContainer()
	itemMenu = fyne.NewContainer()

	testMenu = fyne.NewContainer()
	appIcon,_ = fyne.LoadResourceFromPath("Assets/icon.png")
	profitGraph = canvas.NewImageFromFile("Assets/graph.png")
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	a.SetIcon(appIcon)
	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")
	//a.SetIcon(icon)

	title := widget.NewLabel("Welcome!")
	title.Alignment = fyne.TextAlign(1)

	mainMenu = container.NewVBox(
		title,
		profitGraph,
		canvas.NewImageFromFile("Assets"),
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
		widget.NewLabelWithStyle("Item Menu", fyne.TextAlign(1), fyne.TextStyle{Italic: true}),
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
	//testItemForm := dialog.NewForm("New Item", "Done", "Cancel", []*widget.FormItem ,confirmCallback(), w)
	testMenu = container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Misc", container.NewVBox(
				widget.NewButton("Back", func() {
					w.SetContent(mainMenu)
				}),
				widget.NewButton("Time", func() {
					a.SendNotification(fyne.NewNotification(Data.ConvertDate(time.Now()), Data.ConvertClock(time.Now())))
				}),
				widget.NewButton("Test Notifcation", func() {
					a.SendNotification(fyne.NewNotification("Tree", "I am the lorax, I speak for the tress."))
				}),
				widget.NewButton("Run Test Main", func() {
					Data.TestMain()
				}),
				widget.NewCard("Homies", "You Thought...", widget.NewEntry()),
			)),

			//Shop still not completely
			container.NewTabItem("Shop", container.NewVBox(
				widget.NewLabel("Shopping"),
				//Put code for a binded cart total
				testTitle,
				//Put code for a binded list
				widget.NewButton("New Cart Cart", func() {
					Data.ClearCart(ShoppingCart)
				}),
				widget.NewButton("Show Cart Contents", func() {
					ShoppingCart = Data.AddToCart(13000006057, ShoppingCart)
					fmt.Println(ShoppingCart)
				}),
				widget.NewButton("Add B1 To Cart", func() {
					file, _ := os.Open(Cam.Path + "Online Test 01.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					conID, _ := strconv.Atoi(id)
					ShoppingCart = Data.AddToCart(conID, ShoppingCart)
					total := fmt.Sprint(Data.GetCartTotal(ShoppingCart))
					fmt.Println(ShoppingCart, total)
					testTitle.SetText(id)
				}),
			)),

			//Shop still not completely
			container.NewTabItem("Shop 2", container.NewVSplit(
				container.NewVBox(
					container.NewVBox(
						widget.NewLabel("Shopping")),
				),
				container.NewVBox(
					),
				//Put code for a binded cart total
			)),

			container.NewTabItem("Barcodes", container.NewVBox(
				testTitle,
				widget.NewButton("Camera", func() {
					//Cam.OpenCam()
				}),
				widget.NewButton("Barcode 01", func() {
					file, _ := os.Open(Cam.Path + "Online Test 01.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					newId, _ := strconv.Atoi(id)

					CreateNewItem(newId, w)
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 02", func() {
					file, _ := os.Open(Cam.Path + "Online Test 02.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					newId, _ := strconv.Atoi(id)

					CreateNewItem(newId, w)
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 03", func() {
					file, _ := os.Open(Cam.Path + "Online Test 03.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					newId, _ := strconv.Atoi(id)

					CreateNewItem(newId, w)
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 04", func() {
					file, _ := os.Open(Cam.Path + "Test01.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					newId, _ := strconv.Atoi(id)

					CreateNewItem(newId, w)
					testTitle.SetText("ID: " + id)
				}),
				widget.NewButton("Barcode 05", func() {
					file, _ := os.Open(Cam.Path + "Online Test 05.png")
					img, _, _ := image.Decode(file)
					id := Cam.ReadImage(img).String()
					newId, _ := strconv.Atoi(id)

					CreateNewItem(newId, w)
					testTitle.SetText("ID: " + id)
				}),
			)),

			container.NewTabItem("Info", container.NewVBox(
				widget.NewLabel("Editing Code Info"),
				container.NewHSplit(
					container.NewVScroll(
						container.NewVBox(
							widget.NewButton("Barcode 05", func() {
								//Get the Index of the barcode in the data
								//Fill the menu's placeholders with data from the original id
								file, _ := os.Open(Cam.Path + "Online Test 05.png")
								img, _, _ := image.Decode(file)
								id := Cam.ReadImage(img).String()
								conID, _ := strconv.Atoi(id)

								ModifyItem(conID, w)
								//Grab and display the data from the cells in that row
							}),
						)),
					widget.NewForm(
						widget.NewFormItem("PlaceHolder", widget.NewLabel("Probably going to replace the data displays binded to the data")),
						widget.NewFormItem("Id", widget.NewLabel("ID")),
						widget.NewFormItem("Name", widget.NewEntry()),
						widget.NewFormItem("Price", UI.NewNumEntry()),
						widget.NewFormItem("Cost", UI.NewNumEntry()),
						widget.NewFormItem("Inventory", UI.NewNumEntry()),
					)),
			)),

			container.NewTabItem("Stats", container.NewVBox(
				widget.NewAccordion(widget.NewAccordionItem("Today", widget.NewButton("Today", func() {
					a.SendNotification(fyne.NewNotification("Hello", ""))
				}))),
			)),

		),
	)

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func CreateNewItem(id int, w fyne.Window){
	//password := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	idLabel := widget.NewLabel(strconv.Itoa(id))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("The Name of the Product")
	nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

	priceEntry := UI.NewNumEntry()
	priceEntry.SetPlaceHolder("The Price it will be sold.")

	costEntry := UI.NewNumEntry()
	costEntry.SetPlaceHolder("The Cost of buying this item.")

	inventoryEntry := UI.NewNumEntry()
	inventoryEntry.SetPlaceHolder("The Amount currently in inventory.")

	items := []*widget.FormItem{
		widget.NewFormItem("ID", idLabel),
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Price", priceEntry),
		widget.NewFormItem("Cost", costEntry),
		widget.NewFormItem("Inventory", inventoryEntry),
	}

	dialog.ShowForm("New Item", "Add", "Cancel", items, func(b bool) {
		if !b {
			return
		}

		//log.Println("Please Check the Price, cost, or the amount in inventory")
		fmt.Println("Name, Price, Cost and Inventory have all been Authenticated...")
		fmt.Println("Adding to the database")
		//Call a sort of save method that will take the data and save it to the item data
		//Data.SaveNewItem(id, Data.NewSale(0, "Blue balls", 5.5, 1.25, 2))
		price, cost, inventory := Data.ConvertStringToSale(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
		Data.UpdateData(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Items", 2)
		Data.ReadVal("Items")
		Data.SaveFile()
	}, w)
}

func ModifyItem(id int, w fyne.Window){
	//password := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	idLabel := widget.NewLabel(strconv.Itoa(id))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("The Name of the Product")
	nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

	priceEntry := UI.NewNumEntry()
	priceEntry.SetPlaceHolder("The Price it will be sold.")

	costEntry := UI.NewNumEntry()
	costEntry.SetPlaceHolder("The Cost of buying this item.")

	inventoryEntry := UI.NewNumEntry()
	inventoryEntry.SetPlaceHolder("The Amount currently in inventory.")

	items := []*widget.FormItem{
		widget.NewFormItem("ID", idLabel),
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Price", priceEntry),
		widget.NewFormItem("Cost", costEntry),
		widget.NewFormItem("Inventory", inventoryEntry),
	}

	dialog.ShowForm("Modify Item", "Change", "Cancel", items, func(b bool) {
		if !b {
			fmt.Println("Doesn't work.")
			return
		}

		//log.Println("Please Check the Price, cost, or the amount in inventory")
		fmt.Println("Name, Price, Cost and Inventory have all been Authenticated...")
		fmt.Println("Adding to the database")
		//Call a sort of save method that will take the data and save it to the item data
		//Data.SaveNewItem(id, Data.NewSale(0, "Blue balls", 5.5, 1.25, 2))
		price, cost, inventory := Data.ConvertStringToSale(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
		//Data.UpdateData(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Items", 2)
		Data.ModifyItem(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Items")
		Data.ReadVal("Items")
		Data.SaveFile()
	}, w)
}