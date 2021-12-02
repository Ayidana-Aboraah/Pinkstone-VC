package main

import (
	"business.go/Cam"
	"business.go/Data"
	"business.go/Graph"
	"business.go/UI"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "image/png"
	"net/url"
	"strconv"
)

var (
	mainMenu = container.NewWithoutLayout()
	appIcon,_ = fyne.LoadResourceFromPath("Assets/icon02.png")
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	a.SetIcon(appIcon)
	go Graph.StartServer()

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

	button := widget.NewButton("New Item", func() {
		//Get ID and Convert
		id := Cam.OpenCam()
		conID, _ := strconv.Atoi(id)

		raw := Data.GetData("Items", conID)
		dialog.ShowConfirm("Check (Move middle bar)", "Is this the right item: " + raw[0], func(b bool) {
			if !b{
				return
			}
			//Append the item to the cartList
			cart, _ := cartList.Get()
			cartList.Set(Data.AddToCart(conID, cart))
			bondTotal.Set(Data.GetCartTotal(cart))
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
			quantity := fmt.Sprint(i.Quantity)
			text.SetText(i.Name + " x" + quantity)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				cart, _ := cartList.Get()
				cartList.Set(Data.DecreaseFromCart(val.ID,cart))
			}
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
	listData := Data.GetAllData("Items")
	boundData := binding.BindSaleList(&listData)
	list := widget.NewListWithData(boundData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, widget.NewButton("i", nil), widget.NewLabel(""))
	},
	func(item binding.DataItem, obj fyne.CanvasObject) {
		f := item.(binding.Sale)
		val, _ := f.Get()

		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(val.Name)

		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		btn.OnTapped = func() {
			vals := Data.ConvertSaleToString(val.Price, val.Cost, val.Quantity)
			idLabel.SetText(strconv.Itoa(val.ID))
			nameLabel.SetText(val.Name)
			priceLabel.SetText(vals[0])
			costLabel.SetText(vals[1])
			inventoryLabel.SetText(vals[2])
		}
	})

	box := container.NewVBox(
		title,
		container.NewHSplit(
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
					list,
				),
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
	testLink := widget.NewHyperlink("Profits Graph", u)

	selectionEntry := UI.NewNumEntry()
	selectionEntry.SetPlaceHolder("YYYY/MM/Day")

	scroll := container.NewVScroll(
		container.NewAppTabs(container.NewTabItem("Graphs",
		container.NewVBox(
		widget.NewCard("Item Popularity Chart", "See the ", container.NewVBox(
			selectionEntry,
			widget.NewButton("Graph", func() {
				labels, profits := Data.GetTotalProfit(selectionEntry.Text, 2)
				//colors := []string{"Red", "Blue", "Green", "Purple", "Violet", "Orange", "Indigo", "White", "Black"}
				cats := []string{selectionEntry.Text}
				fmt.Println(labels, profits)

				Graph.Labels = &cats
				Graph.Categories = &labels
				Graph.Inputs = &profits
			}),
			//Put a graph here
			testLink,
			)),
		),
	),
	container.NewTabItem("Numbers",
		container.NewVBox(
			)),
		))

	return scroll
}