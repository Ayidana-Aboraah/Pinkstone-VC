package main

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Database"
	"BronzeHermes/Graph"
	"BronzeHermes/UI"
	"fmt"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("Bronze Hermes")

	UI.HandleErrorWithMessage(Database.LoadData(), "Failed to Load Data", a.NewWindow("Error MSG"))

	go Graph.StartServer()

	fmt.Println(Database.Databases)

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")
	w.SetOnClosed(Graph.StopSever)

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", makeMainMenu(a)),
		container.NewTabItem("Shop", makeShoppingMenu(w)),
		container.NewTabItem("Inventory", makeInfoMenu(w)),
		container.NewTabItem("Statistics", makeStatsMenu()),
	)))

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Save Backup Data", func() {
			go func() {
				err := Database.BackUpAllData()
				UI.HandleErrorWithMessage(err, "Failed to Load Data", a.NewWindow("Error MSG"))
			}()
		}),
		widget.NewButton("Load Backup Data", func() {
			err := Database.LoadBackUp()
			UI.HandleErrorWithMessage(err, "Failed to Load Data", a.NewWindow("Error MSG"))
		}),
		widget.NewButton("Display Database", func() { fmt.Println(Database.Databases) }),
		widget.NewButton("Quit", a.Quit))
}

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject {
	var shoppingCart []Database.Sale
	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	cartList := binding.BindSaleList(&shoppingCart)

	list := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil), widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {})

	list.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		title.SetText(fmt.Sprintf("Cart Total: %f", Database.GetCartTotal(shoppingCart)))
		cartList.Reload()
		list.Unselect(id)
	}

	list.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		val := shoppingCart[id]

		text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
		btn.OnTapped = func() {
			shoppingCart = Database.DecreaseFromCart(val, shoppingCart)
			title.SetText(fmt.Sprintf("Cart Total: %f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
			cartList.Reload()
			list.Refresh()
		}
	}

	screen := container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(list),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b {
						return
					}
					shoppingCart = Database.BuyCart(shoppingCart)
					cartList.Reload()
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func() {
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
				cartList.Reload()
			}),
			widget.NewButton("New Item", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				item := Database.FindItem(uint64(id))

				dialog.ShowCustomConfirm("Check", "Yes", "No", container.NewVBox(widget.NewLabel("Is this the right item: "+Database.NameKeys[item.ID])),
					func(b bool) {
						if !b {
							return
						}

						shoppingCart = Database.AddToCart(item, shoppingCart)
						title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
						cartList.Reload()
						list.Refresh()
					}, w)
			}),
		),
	)
	return screen
}

func makeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	inventoryLabel := widget.NewLabel("Inventory")

	title := widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true})

	boundData := binding.BindSaleList(&Database.Databases[0])
	inventoryList := widget.NewListWithData(boundData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"))
	},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Sale)
			val, _ := f.Get()
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Database.NameKeys[val.ID])
		})

	inventoryList.OnSelected = func(id widget.ListItemID) {
		item := Database.Databases[0][id]
		values := Database.ConvertSale(item)

		idLabel.SetText(strconv.Itoa(int(item.ID)))
		nameLabel.SetText(Database.NameKeys[item.ID])
		priceLabel.SetText(values[0])
		costLabel.SetText(values[1])
		inventoryLabel.SetText(values[2])
	}

	box := container.New(layout.NewGridLayout(2),
		container.NewVBox(
			title,
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			inventoryLabel,
			widget.NewButton("Modify", func() {
				//Creates item menu
				func() {
					conID, _ := strconv.Atoi(idLabel.Text)

					idLabel := widget.NewLabel(strconv.Itoa(conID))

					nameEntry := widget.NewEntry()
					nameEntry.SetPlaceHolder("Product Name with _ for spaces.")
					nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

					priceEntry := UI.NewNumEntry("The Price it will be sold.")
					costEntry := UI.NewNumEntry("The Cost of buying this item.")
					inventoryEntry := UI.NewNumEntry("The Amount currently in inventory.")

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

						price, cost, inventory := Database.ConvertString(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
						newItem := Database.Sale{ID: uint64(conID), Price: price, Cost: cost, Quantity: inventory}

						Database.Databases[0] = append(Database.Databases[0], newItem)
						Database.Databases[2] = append(Database.Databases[2], newItem)
						Database.AddKey(uint64(conID), nameEntry.Text)

						boundData.Set(Database.Databases[0])
						inventoryList.Refresh()

						Database.SaveData()
						fmt.Println("Saved to Database.")
					}, w)
				}()
			}),
			widget.NewButton("Camera", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				result := Database.FindItem(uint64(id))
				res := Database.ConvertSale(result)

				idLabel.SetText(strconv.Itoa(id))
				nameLabel.SetText(Database.NameKeys[result.ID])
				priceLabel.SetText(res[0])
				costLabel.SetText(res[1])
				inventoryLabel.SetText(res[2])
			}),
		),
		container.NewMax(
			inventoryList,
		))
	return box
}

func makeStatsMenu() fyne.CanvasObject {
	u, err := url.Parse("http://localhost:8081/line")
	UI.HandleError(err)
	r, err := url.Parse("http://localhost:8081/pie")
	UI.HandleError(err)

	link := widget.NewHyperlink("Go To Graph", u)

	selectionEntry := UI.NewNumEntry("Year/Month")

	var lineDataSelectType int
	dataSelectOptions := widget.NewSelect([]string{"Revenue", "Cost", "Profit"}, func(dataType string) {
		switch dataType {
		case "Revenue":
			lineDataSelectType = 0
		case "Cost":
			lineDataSelectType = 1
		case "Profit":
			lineDataSelectType = 2
		}
	})

	days := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
		"15", "16", "17", "18", "19", "20", "21", "22", "23", "24",
		"25", "26", "27", "28", "29", "30", "31"}

	var buttonType int

	scroll := container.NewMax(container.NewVBox(
		widget.NewCard("Data Over Time", "", container.NewVBox(
			selectionEntry,
			widget.NewSelect([]string{"Items Graph", "Price Changes", "Item Popularity", "Item Sales"}, func(graph string) {
				switch graph {
				case "Items Graph":
					buttonType = 0
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Price Changes":
					buttonType = 1
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Item Popularity":
					buttonType = 2
					link.URL = r
					dataSelectOptions.Hidden = false
				case "Item Sales":
					buttonType = 3
					link.URL = r
					dataSelectOptions.Hidden = true
				case "Sales Over Time":
					buttonType = 4
					link.URL = u
					dataSelectOptions.Hidden = true
				}
			}),
			dataSelectOptions,
			widget.NewButton("Graph", func() {
				switch buttonType {
				case 0:
					labels, results := Database.GetLine(selectionEntry.Text, lineDataSelectType, Database.Databases[1])

					Graph.Labels = days
					Graph.Categories = labels
					Graph.LineInputs = results
				case 1:
					labels, results := Database.GetLine(selectionEntry.Text, lineDataSelectType, Database.Databases[2])

					Graph.Labels = days
					Graph.Categories = labels
					Graph.LineInputs = results
				case 2:
					labels, profits := Database.GetPricePie(selectionEntry.Text, lineDataSelectType)

					Graph.Labels = labels
					Graph.Inputs = profits
				case 3:
					labels, sales := Database.GetSalesPie(selectionEntry.Text)

					Graph.Labels = labels
					Graph.Inputs = sales
				case 4:
					labels, sales := Database.GetSalesLine(selectionEntry.Text)

					Graph.Labels = labels
					Graph.LineInputs = sales
				}
			}),
			link,
		)),
	))

	return container.NewVScroll(scroll)
}
