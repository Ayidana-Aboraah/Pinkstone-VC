package main

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Data"
	"BronzeHermes/Graph"
	"BronzeHermes/UI"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/360EntSecGroup-Skylar/excelize"
	"net/url"
	"strconv"
)

var (
	appIcon, _ = fyne.LoadResourceFromPath("Assets/icon02.png")
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	a.SetIcon(appIcon)
	go Graph.StartServers()

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")

	if Data.Err != nil {
		Data.SaveBackUp("BackupAppData.xlsx", "AppData.xlsx")
		Data.F, Data.Err = excelize.OpenFile("Assets/AppData.xlsx")
		dialog.ShowError(Data.Err, w)
	}

	mainMenu := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Main", makeMainMenu(a)),

			container.NewTabItem("Shop", makeShoppingMenu(w)),

			container.NewTabItem("Inventory", makeInfoMenu(w)),

			container.NewTabItem("Statistics", makeStatsMenu()),
		))

	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject {
	box := container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Back Up App Data", func() {
			go Data.SaveBackUp("AppData.xlsx", "BackupAppData.xlsx")
		}),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	)
	return box
}

func createItemMenu(id int, w fyne.Window, boundData binding.ExternalSaleList, list *widget.List) {
	idLabel := widget.NewLabel(strconv.Itoa(id))

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

		fmt.Println("Name, Price, Cost and Inventory have all been Authenticated...")
		fmt.Println("Adding to the database")

		price, cost, inventory := Data.ConvertStringToSale(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
		Data.UpdateData(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Items", 2)
		Data.UpdateData(Data.NewSale(id, nameEntry.Text, price, cost, inventory), "Price Log", 0)
		boundData.Set(Data.GetAllData("Items", 0))
		list.Refresh()

		Data.ReadVal("Items")
		Data.SaveFile()
	}, w)
}

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject {
	var shoppingCart []Data.Sale

	cartList := binding.BindSaleList(&shoppingCart)

	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	list := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil),
				widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {})

	list.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		title.SetText(fmt.Sprintf("Cart Total: %f", Data.GetCartTotal(shoppingCart)))
		cartList.Reload()
		list.Unselect(id)
	}

	list.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		val := shoppingCart[id]

		text.SetText(val.Name + " x" + strconv.Itoa(val.Quantity))
		btn.OnTapped = func() {
			shoppingCart = Data.DecreaseFromCart(val, shoppingCart)
			title.SetText(fmt.Sprintf("Cart Total: %f", Data.GetCartTotal(shoppingCart)))
			text.SetText(val.Name + " x" + strconv.Itoa(val.Quantity))
			cartList.Reload()
			list.Refresh()
		}
	}

	button := widget.NewButton("New Item", func() {
		//Get ID and Convert
		id, err, msg := Cam.OpenCam()
		if err != nil{
			dialog.ShowError(err, w)
			dialog.ShowInformation("Error: 01", "Camera Issue is present: " + msg, w)
		}

		stringID := id.String()

		conID, _ := strconv.Atoi(stringID)

		raw := Data.GetAllData("Items", conID)
		priceEntry := UI.NewNumEntry(fmt.Sprint(raw[0].Price))
		priceEntry.Text = fmt.Sprint(raw[0].Price)

		dialog.ShowCustomConfirm("Check (Move middle bar)", "Yes", "No",
			container.NewVBox(
				widget.NewLabel("Is this the right item: "+raw[0].Name),
				priceEntry,
			),
			func(b bool) {
				if !b {
					return
				}
				//Append the item to the cartList
				newPrice, _, _ := Data.ConvertStringToSale(priceEntry.Text, "", "")
				raw[0].Price = newPrice
				raw[0].Quantity = 1

				shoppingCart = Data.AddToCart(raw[0], shoppingCart)
				title.SetText(fmt.Sprintf("Cart Total: %1.10f", Data.GetCartTotal(shoppingCart)))
				cartList.Reload()
				list.Refresh()
			}, w)
	})

	screen := container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(list),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b {
						return
					}
					shoppingCart = Data.BuyCart(shoppingCart)
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", Data.GetCartTotal(shoppingCart)))
					cartList.Reload()
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func() {
				shoppingCart = Data.ClearCart(shoppingCart)
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Data.GetCartTotal(shoppingCart)))
				cartList.Reload()
			}),
			button,
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

	inventoryData := Data.GetAllData("Items", 0)
	boundData := binding.BindSaleList(&inventoryData)
	inventoryList := widget.NewListWithData(boundData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"))
	},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Sale)
			val, _ := f.Get()
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(val.Name)
		})

	inventoryList.OnSelected = func(id widget.ListItemID) {
		item := inventoryData[id]
		values := Data.ConvertSaleToString(item.Price, item.Cost, item.Quantity)

		idLabel.SetText(strconv.Itoa(item.ID))
		nameLabel.SetText(item.Name)
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
				conID, _ := strconv.Atoi(idLabel.Text)
				createItemMenu(conID, w, boundData, inventoryList)
			}),
			widget.NewButton("Camera", func() {
				id, err, msg := Cam.OpenCam()
				if err != nil{
					dialog.ShowError(err, w)
					dialog.ShowInformation("Error: 01", "Camera Issue is present: " + msg, w)
				}

				stringID := id.String()
				conID, _ := strconv.Atoi(stringID)

				results := Data.GetAllData("Items", conID)
				res := Data.ConvertSaleToString(results[0].Price, results[0].Cost, results[0].Quantity)

				idLabel.SetText(stringID)
				nameLabel.SetText(results[0].Name)
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
	u, _ := url.Parse("http://localhost:8081/line")
	r, _ := url.Parse("http://localhost:8081/pie")

	lineLink := widget.NewHyperlink("Profits Graph", u)
	pieLink := widget.NewHyperlink("Total Sales graph", r)

	lineSelectionEntry := UI.NewNumEntry("Year/Month")
	pieSelectionEntry := UI.NewNumEntry("YYYY/MM/Day")

	var lineDataSelectType int
	dataSelectOptions := widget.NewSelect([]string{"Revenue", "Cost", "Profit"},func(dataType string){
		if dataType == "Revenue"{lineDataSelectType = 0}
		if dataType == "Cost"{lineDataSelectType = 1}
		if dataType == "Profit"{lineDataSelectType = 2}
	})

	totalRevLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	totalCostLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	totalProfitLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	scroll := container.NewMax(
			container.NewVBox(
				widget.NewCard("Items Graph", "", container.NewVBox(
					lineSelectionEntry,
					dataSelectOptions,
					widget.NewButton("Graph", func() {
						results, labels := Data.GetProfitForTimes(lineDataSelectType, "Report Data", lineSelectionEntry.Text)
						days := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
							"15", "16", "17", "18", "19", "20", "21", "22", "23", "24",
							"25", "26", "27", "28", "29", "30", "31"}

						Graph.Labels = &days
						Graph.Categories = labels
						Graph.LineInputs = &results
					}),
					//Put a graph here
					lineLink,
				)),

				widget.NewCard("Price Changes", "", container.NewVBox(
					lineSelectionEntry,
					dataSelectOptions,
					widget.NewButton("Graph", func() {
						results, labels := Data.GetProfitForTimes(lineDataSelectType, "Price Log", lineSelectionEntry.Text)
						days := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
							"15", "16", "17", "18", "19", "20", "21", "22", "23", "24",
							"25", "26", "27", "28", "29", "30", "31"}

						Graph.Labels = &days
						Graph.Categories = labels
						Graph.LineInputs = &results
					}),
					lineLink,
				)),

				widget.NewCard("Item Popularity", "", container.NewVBox(
					pieSelectionEntry,
					dataSelectOptions,
					widget.NewButton("Graph", func() {
						profits, labels := Data.GetAllProfits(pieSelectionEntry.Text)

						Graph.Labels = &labels
						Graph.Inputs = &profits[lineDataSelectType]
					}),
					pieLink,
				)),
				widget.NewCard("Totals", "", container.NewVBox(
					pieSelectionEntry,
					widget.NewButton("Graph", func() {
						data, _ := Data.GetAllProfits(pieSelectionEntry.Text)

						revenue := fmt.Sprint(data[0])
						cost := fmt.Sprint(data[1])
						profit := fmt.Sprint(data[2])

						totalProfitLabel.SetText("Total Profit: " + profit)
						totalRevLabel.SetText("Total Revenue: " + revenue)
						totalCostLabel.SetText("Total Cost: " + cost)
					}),
					totalRevLabel,
					totalCostLabel,
					totalProfitLabel,
				)),
			))

	return container.NewVScroll(scroll)
}