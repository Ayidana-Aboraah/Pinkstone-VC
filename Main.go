package main

import (
	"business.go/Cam"
	"business.go/Data"
	"business.go/UI"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
)

var (
	mainMenu = container.NewWithoutLayout()
	appIcon,_ = fyne.LoadResourceFromPath("Assets/icon02.png")
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	a.SetIcon(appIcon)

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")

	if Data.Err != nil{
		fmt.Println(Data.Err)
		dialog.ShowError(Data.Err, w)
	}

	mainMenu = container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Main", makeMainMenu(a)),
			//Shop still not completely
			container.NewTabItem("Shop", makeShoppingMenu(w)),

			container.NewTabItem("Info", makeInfoMenu(w)),

			container.NewTabItem("Stats", makeStatsMenu(w)),

		))

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject{
	box := container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Back Up App Data", func() {
			//Don't forget to change the source file name when switching from test file to normal file
			Data.SaveBackUp("TestAppData.xlsx", "BackupAppData.xlsx")
		}),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	)
	return box
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
	bondTotal := binding.BindFloat(&total)

	cartList := binding.BindSaleList(&[]Data.Sale{})

	title := widget.NewLabelWithData(binding.FloatToString(bondTotal))

	list := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil),
				widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Sale)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			i, _ := f.Get()
			quantity := fmt.Sprint(i.Quantity)
			text.SetText(i.Name + " x" + quantity)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				cart, _ := cartList.Get()
				cartList.Set(Data.DecreaseFromCart(val.ID,cart))
				//ShoppingCart = Data.DecreaseFromCart(val.ID,*ShoppingCart)
			}
		})

/*
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
				ShoppingCart = Data.DecreaseFromCart(f.ID, *ShoppingCart)
				text.Refresh()
			}
		})
 */

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
			cart, _ := cartList.Get()
			cartList.Set(Data.AddToCart(conID, cart))
			bondTotal.Set(Data.GetCartTotal(cart))
		},w)
	})

	split := container.NewVSplit(
			container.NewGridWithColumns(1,
				title,
				list,
			),

		container.NewHBox(
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b{
						fmt.Println("Understandable.")
						return
					}
					cart, _ := cartList.Get()
					cartList.Set(Data.BuyCart(cart))
					bondTotal.Set(Data.GetCartTotal(cart))
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func(){
				cart, _ := cartList.Get()
				cartList.Set(Data.ClearCart(cart))
				bondTotal.Set(Data.GetCartTotal(cart))
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
func makeStatsMenu(w fyne.Window) fyne.CanvasObject {
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