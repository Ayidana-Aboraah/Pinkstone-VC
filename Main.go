package main

import (
	"business.go/Cam"
	"business.go/Data"
	"business.go/UI"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/png"
	"net/url"
	"os"
	"strconv"
)

var ShoppingCart []*Data.Sale

var (
	mainMenu = fyne.NewContainer()
	dataMenu = fyne.NewContainer()
	shopMenu = fyne.NewContainer()
	settingsMenu = fyne.NewContainer()

	testMenu = fyne.NewContainer()
	appIcon,_ = fyne.LoadResourceFromPath("Assets/icon02.png")
	//profitGraph = canvas.NewImageFromFile("Assets/profitGraph.png")
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
		//Put a link to the profit web server
		widget.NewButton("Shopping", func() {
			w.SetContent(shopMenu)
		}),
		widget.NewButton("Statistics", func() {
			w.SetContent(dataMenu)
		}),
		widget.NewButton("Settings", func() {
			w.SetContent(settingsMenu)
		}),
		widget.NewButton("Test", func() {
			w.SetContent(testMenu)
		}),
		widget.NewButton("Quit", func() {
			w.Close()
		}),
	)

	dataMenu = container.NewVBox(
		makeStatsMenu(a, w),
	)

	settingsMenu = container.NewVBox(
		makeInfoMenu(w),
	)

	testMenu = container.NewVBox(
		container.NewAppTabs(
			//Shop still not completely
			container.NewTabItem("Shop", makeShoppingMenu(w)),

			container.NewTabItem("Info", makeInfoMenu(w)),

			container.NewTabItem("Stats", makeStatsMenu(a, w)),

		),
	)

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func createItemMenu(id int, w fyne.Window){
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

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject{
	total := 0.0
	conTotal := fmt.Sprint(total)
	title := widget.NewLabel("Cart Total: " + conTotal)

	list := widget.NewList(func() int {return len(ShoppingCart)},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil),
				widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			f := ShoppingCart[id]
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)

			quantity := fmt.Sprint(f.Quantity)
			text.SetText(f.Name + " x" + quantity)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				ShoppingCart = Data.DecreaseFromCart(f.ID, ShoppingCart)
				text.Refresh()
			}
		})

	button := widget.NewButton("New Item", func() {
		//Get ID and Convert
		id := Cam.OpenCam()
		conID, _ := strconv.Atoi(id)

		raw := Data.GetData("Items", conID)
		dialog.ShowConfirm("Check (Move the bar in the middle to update the list)", "Is this the right item: " + raw[0], func(b bool) {
			if !b{
				return
			}
			//Append the item to the cartList
			ShoppingCart = Data.AddToCart(conID, ShoppingCart)
		},w)
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
			}),
			button,
		),
	)
	return split
}

func makeInfoMenu(w fyne.Window) fyne.CanvasObject{
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	inventoryLabel := widget.NewLabel("Inventory")

	title := widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1),fyne.TextStyle{Bold: true})
	//Create a list of all registered items

	box := container.NewVBox(
		title,
		container.NewHSplit(
			container.NewVScroll(
				container.NewVBox(
					widget.NewButton("Camera", func() {
						id := Cam.OpenCam()
						conID, _ := strconv.Atoi(id)

						res := Data.GetData("Items", conID)

						idLabel.SetText(id)
						nameLabel.SetText(res[0])
						priceLabel.SetText(res[1])
						costLabel.SetText(res[2])
						inventoryLabel.SetText(res[3])
					}),
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
					createItemMenu(conID, w)
				}),
			)),
	)
	return  box
}

//Finish setting up graph stuff for it
func makeStatsMenu(a fyne.App, w fyne.Window) fyne.CanvasObject {
	u, _ := url.Parse("http://localhost:8081")
	testLink := widget.NewHyperlink("Random Line Graph", u)

	back := widget.NewButton("Back", func() {
		w.SetContent(mainMenu)
	})

	selectionEntry := UI.NewNumEntry()

	box := container.NewVBox(
		selectionEntry,
		widget.NewButton("New Graph", func() {
			//rev, cos, prof := Data.GetTotalProfit(selectionEntry.Text)
			go UI.StartServer()
		}),
		//Put a graph here
		testLink,
		back,
	)

	return box
}