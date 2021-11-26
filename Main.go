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
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/png"
	"net/url"
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
	appIcon,_ = fyne.LoadResourceFromPath("Assets/icon02.png")
	profitGraph = canvas.NewImageFromFile("Assets/graph.png")
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	a.SetIcon(appIcon)

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")

	title := widget.NewLabel("Welcome!")
	title.Alignment = fyne.TextAlign(1)

	mainMenu = container.NewVBox(
		title,
		profitGraph,
		canvas.NewImageFromFile("Assets"),
		widget.NewButton("Shopping", func() {
			w.SetContent(shopMenu)
		}),
		widget.NewButton("Statistics", func() {
			w.SetContent(dataMenu)
		}),
		widget.NewButton("Settings", func() {
			//w.SetContent(settingsMenu)
		}),
		widget.NewButton("Test", func() {
			w.SetContent(testMenu)
		}),
		widget.NewButton("Camera", func() {
			Cam.StartCamera()
		}),
		widget.NewButton("Quit", func() {
			w.Close()
		}),
	)

	dataMenu = container.NewVBox(
		widget.NewLabelWithStyle("Data", 1, fyne.TextStyle{}),
		widget.NewButton("Back", func() {
			title.SetText("Bronze Hermes")
			w.SetContent(mainMenu)
		}),
	)

	itemMenu = container.NewVBox(
	)

	testMenu = container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Misc", container.NewVBox(
				widget.NewButton("Back", func() {
					w.SetContent(mainMenu)
				}),
				widget.NewButton("Time", func() {
					a.SendNotification(fyne.NewNotification(Data.ConvertDate(time.Now()), Data.ConvertClock(time.Now())))
				}),
				widget.NewButton("Notification", func() {
					a.SendNotification(fyne.NewNotification("Tree", "I am the lorax, I speak for the tress."))
				}),
				widget.NewButton("Run Test Main", func() {
					Data.TestMain()
				}),
				widget.NewCard("Trash Afton", "You wish", UI.NewNumEntry()),
				widget.NewButton("Start Graph", func() {
					go UI.StartServer()
				}),
				widget.NewButton("Open Graph", func() {
					path := url.URL{Path: "http://localhost:8081/"}
					a.OpenURL(&path)
				}),
			)),

			//Shop still not completely
			container.NewTabItem("Shop 2", makeShoppingMenu(w)),

			container.NewTabItem("Barcodes", makeBarcodeMenu(w)),

			container.NewTabItem("Info", makeInfoMenu(w)),

			container.NewTabItem("Stats", makeStatsMenu()),

		),
	)

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func CreateNewItem(id int, w fyne.Window){
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

		fmt.Println("Name, Price, Cost and Inventory have all been Authenticated...")
		fmt.Println("Adding to the database")

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
			fmt.Println("Item could not be modified, check if it exists.")
			return
		}

		fmt.Println("Name, Price, Cost and Inventory have all been Authenticated...")
		fmt.Println("Adding to the database")

		price, cost, inventory := Data.ConvertStringToSale(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
		Data.ModifyItem(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Items")

		Data.ReadVal("Items")
		Data.SaveFile()
	}, w)
}

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject{
	total := 0.0
	conTotal := fmt.Sprint(total)
	title := widget.NewLabel("Cart Total: " + conTotal)
	cartList := binding.BindSaleList(&[]Data.Sale{})

	button := widget.NewButton("New Item", func() {
		//Get ID and Convert
		id := Cam.OpenCam()
		conID, _ := strconv.Atoi(id)

		raw := Data.GetData("Items", conID)
		dialog.ShowConfirm("Check", "Is this the right item: " + raw[0], func(b bool) {
			if !b{
				return
			}

			//Scan Item
			price, cost, inventory := Data.ConvertStringToSale(raw[1], raw[2], raw[3])

			//Append the item to the cartList
			red := Data.Sale{ID: conID, Name: raw[0], Price: price, Cost: cost, Quantity: inventory}
			cartList.Append(red)
			ShoppingCart = append(ShoppingCart, &red)
		},w)
	})

	list := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil),
				widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Sale)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			i, _ := f.Get()
			text.SetText(i.Name)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				ShoppingCart = Data.DecreaseFromCart(val.ID, ShoppingCart)
				fmt.Println(val, ShoppingCart)
			}
		})

	split := container.NewVSplit(
		container.NewVScroll(
			//Put code for a binded cart total
			container.NewGridWithColumns(1,
				title,
				list,
			),
		),
		container.NewHBox(
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b{
						fmt.Println("Understandable.")
						return
					}

					ShoppingCart = Data.BuyCart(ShoppingCart)
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func(){
				ShoppingCart = Data.ClearCart(ShoppingCart)
				//Remove the buttons somehow
			}),
			button,
		),
	)
	return split
}

func makeBarcodeMenu(w fyne.Window) fyne.CanvasObject{
	box := container.NewVBox(
		widget.NewLabel("Barcodes"),
		widget.NewButton("Camera", func() {
			id := Cam.OpenCam()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
		widget.NewButton("Barcode 01", func() {
			file, _ := os.Open(Cam.Path + "Online Test 01.png")
			img, _, _ := image.Decode(file)
			id := Cam.ReadImage(img).String()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
		widget.NewButton("Barcode 02", func() {
			file, _ := os.Open(Cam.Path + "Online Test 02.png")
			img, _, _ := image.Decode(file)
			id := Cam.ReadImage(img).String()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
		widget.NewButton("Barcode 03", func() {
			file, _ := os.Open(Cam.Path + "Online Test 03.png")
			img, _, _ := image.Decode(file)
			id := Cam.ReadImage(img).String()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
		widget.NewButton("Barcode 04", func() {
			file, _ := os.Open(Cam.Path + "Test01.png")
			img, _, _ := image.Decode(file)
			id := Cam.ReadImage(img).String()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
		widget.NewButton("Barcode 05", func() {
			file, _ := os.Open(Cam.Path + "Online Test 05.png")
			img, _, _ := image.Decode(file)
			id := Cam.ReadImage(img).String()
			newId, _ := strconv.Atoi(id)

			CreateNewItem(newId, w)
		}),
	)
	return box
}

func makeInfoMenu(w fyne.Window) fyne.CanvasObject{
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	inventoryLabel := widget.NewLabel("Inventory")

	box := container.NewVBox(
		widget.NewLabel("Editing Code Info"),
		container.NewHSplit(
			container.NewVScroll(
				container.NewVBox(
					widget.NewButton("Barcode 01", func() {
						file, _ := os.Open(Cam.Path + "Online Test 01.png")
						img, _, _ := image.Decode(file)
						id := Cam.ReadImage(img).String()
						conID, _ := strconv.Atoi(id)

						//ModifyItem(conID, w)
						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
					widget.NewButton("Barcode 02", func() {
						file, _ := os.Open(Cam.Path + "Online Test 02.png")
						img, _, _ := image.Decode(file)
						id := Cam.ReadImage(img).String()
						conID, _ := strconv.Atoi(id)

						//ModifyItem(conID, w)
						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
					widget.NewButton("Barcode 03", func() {
						file, _ := os.Open(Cam.Path + "Online Test 03.png")
						img, _, _ := image.Decode(file)
						id := Cam.ReadImage(img).String()
						conID, _ := strconv.Atoi(id)

						//ModifyItem(conID, w)
						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
					widget.NewButton("Barcode 04", func() {
						file, _ := os.Open(Cam.Path + "Test01.png")
						img, _, _ := image.Decode(file)
						id := Cam.ReadImage(img).String()
						conID, _ := strconv.Atoi(id)

						//ModifyItem(conID, w)
						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
					widget.NewButton("Barcode 05", func() {
						file, _ := os.Open(Cam.Path + "Online Test 05.png")
						img, _, _ := image.Decode(file)
						id := Cam.ReadImage(img).String()
						conID, _ := strconv.Atoi(id)

						//ModifyItem(conID, w)
						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
				)),
			container.NewVBox(
				idLabel,
				nameLabel,
				priceLabel,
				costLabel,
				inventoryLabel,
				widget.NewButton("Modify", func() {
					conID, _ := strconv.Atoi(idLabel.Text)
					ModifyItem(conID, w)
				}),
			)),
	)
	return  box
}

//Finish setting up graph stuff for it
func makeStatsMenu() fyne.CanvasObject{
	selectionEntry := UI.NewNumEntry()

	box := container.NewVBox(
		selectionEntry,
		widget.NewButton("New Graph", func() {
			Data.GetTotalProfit(selectionEntry.Text)
		}),
		//Put a graph here
	)
	return box
}