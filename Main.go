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
	"net/url"
	"os"
	"strconv"
)

var (
	mainMenu   = container.NewWithoutLayout()
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
		//Replace TestAppData with normal App data when ready
		os.Remove("TestAppData.xlsx")
		//os.Remove("Assets/AppData.xlsx")
		Data.SaveBackUp("Assets/BackupAppData.xlsx", "TestAppData.xlsx")
		//Data.SaveBackUp("BackupAppData.xlsx", "AppData.xlsx")
		fmt.Println(Data.Err)
		dialog.ShowError(Data.Err, w)
	}

	mainMenu = container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Main", makeMainMenu(a)),
			//Shop still not completely
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
		widget.NewButton("Import Backup", func() {

		}),
		widget.NewButton("Export BackUp", func() {

		}),
		widget.NewButton("Back Up App Data", func() {
			//Don't forget to change the source file name when switching from test file to normal file
			go Data.SaveBackUp("TestAppData.xlsx", "BackupAppData.xlsx")
			//go Data.SaveBackUp("AppData.xlsx", "BackupAppData.xlsx")
		}),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	)
	return box
}

func createItemMenu(id int, w fyne.Window) {
	idLabel := widget.NewLabel(strconv.Itoa(id))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("The Name of the Product")
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

		Data.ReadVal("Items")
		Data.SaveFile()
	}, w)
}

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject {

	cartList := binding.BindSaleList(&[]Data.Sale{})

	title := widget.NewLabelWithStyle("Cart Total: ", fyne.TextAlignCenter, fyne.TextStyle{})

	button := widget.NewButton("New Item", func() {
		//Get ID and Convert
		id := Cam.OpenCam()

		conID, _ := strconv.Atoi(id)

		raw := Data.GetAllData("Items", conID)
		priceEntry := UI.NewNumEntry(fmt.Sprint(raw[0].Price))
		priceEntry.Text = fmt.Sprint(raw[0].Price)

		/*dialog.ShowConfirm("Check (Move middle bar)", "Is this the right item: "+raw[0].Name, func(b bool) {
			if !b {
				return
			}
			//Append the item to the cartList
			cart, _ := cartList.Get()
			cartList.Set(Data.AddToCart(conID, cart))
			cart, _ = cartList.Get()
			title.SetText(fmt.Sprintf("Cart Total: %0.0f",  Data.GetCartTotal(cart)))
		}, w)
		*/
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

				cart, _ := cartList.Get()
				cartList.Set(Data.AddToCart(raw[0], cart))
				cart, _ = cartList.Get()
				title.SetText(fmt.Sprintf("Cart Total: %1.10f", Data.GetCartTotal(cart)))
			}, w)
	})

	list := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil),
				widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s := item.(binding.Sale)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			i, _ := s.Get()

			text.SetText(i.Name + " x" + strconv.Itoa(i.Quantity))

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := s.Get()
				cart, _ := cartList.Get()
				cartList.Set(Data.DecreaseFromCart(val, cart))
				cart, _ = cartList.Get()
				title.SetText(fmt.Sprintf("Cart Total: %f", Data.GetCartTotal(cart)))
				text.SetText(i.Name + " x" + strconv.Itoa(i.Quantity))
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
					if !b {
						return
					}
					cart, _ := cartList.Get()
					cartList.Set(Data.BuyCart(cart))
					cart, _ = cartList.Get()
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", Data.GetCartTotal(cart)))
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func() {
				cart, _ := cartList.Get()
				cartList.Set(Data.ClearCart(cart))
				cart, _ = cartList.Get()
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Data.GetCartTotal(cart)))
			}),
			button,
		),
	)
	return split
}

func makeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	inventoryLabel := widget.NewLabel("Inventory")

	title := widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true})

	//Create a list of all registered items
	listData := Data.GetAllData("Items", 0)

	boundData := binding.BindSaleList(&listData)
	list := widget.NewListWithData(boundData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"))
	},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Sale)
			val, _ := f.Get()
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(val.Name)
		})

	/*
	list := widget.NewList(func() int { return len(listData) },
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel("Name"))
		}, func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(listData[id].Name)
		})
	*/

	list.OnSelected = func(id widget.ListItemID) {
		val := listData[id]
		vals := Data.ConvertSaleToString(val.Price, val.Cost, val.Quantity)

		idLabel.SetText(strconv.Itoa(val.ID))
		nameLabel.SetText(val.Name)
		priceLabel.SetText(vals[0])
		costLabel.SetText(vals[1])
		inventoryLabel.SetText(vals[2])
	}

	box := container.New(layout.NewGridLayoutWithRows(2),
		container.NewVBox(
			title,
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			inventoryLabel,
			widget.NewButton("Modify", func() {
				conID, _ := strconv.Atoi(idLabel.Text)
				createItemMenu(conID, w)
			}),
		),
		container.NewVBox(
			widget.NewButton("Camera", func() {
				id := Cam.OpenCam()
				conID, _ := strconv.Atoi(id)

				results := Data.GetAllData("Items", conID)
				res := Data.ConvertSaleToString(results[0].Price, results[0].Cost, results[0].Quantity)

				idLabel.SetText(id)
				nameLabel.SetText(results[0].Name)
				priceLabel.SetText(res[0])
				costLabel.SetText(res[1])
				inventoryLabel.SetText(res[2])
			}),
			list,
		))

	return box
}

//Finish setting up graph stuff for it
func makeStatsMenu() fyne.CanvasObject {
	u, _ := url.Parse("http://localhost:8081/line")
	r, _ := url.Parse("http://localhost:8081/pie")

	lineLink := widget.NewHyperlink("Profits Graph", u)
	pieLink := widget.NewHyperlink("Pie Graph", r)

	lineSelectionEntry := UI.NewNumEntry("Year/Month")

	pieSelectionEntry := UI.NewNumEntry("YYYY/MM/Day")

	totalRevLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	totalCostLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	totalProfitLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	scroll := container.NewVScroll(
		container.NewAppTabs(container.NewTabItem("Graphs",
			container.NewVBox(
				widget.NewCard("Profit Graph", "", container.NewVBox(
					lineSelectionEntry,
					widget.NewButton("Graph", func() {
						results, labels := Data.GetProfitForTimes(0, "Report Data", lineSelectionEntry.Text)
						days := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
							"15", "16", "17", "18", "19", "20", "21", "22", "23", "24",
							"25", "26", "27", "28", "29", "30", "31"}

						fmt.Println(results)

						Graph.Labels = &days
						Graph.Categories = &labels
						Graph.LineInputs = &results
					}),
					//Put a graph here
					lineLink,
				)),

				widget.NewCard("Price Changes", "", container.NewVBox(
					lineSelectionEntry,
					widget.NewButton("Graph", func() {
					}),
				)),

				widget.NewCard("Item Popularity", "", container.NewVBox(
					pieSelectionEntry,
					widget.NewButton("Graph", func() {
						profits, labels := Data.GetAllProfits(pieSelectionEntry.Text)

						Graph.Labels = &labels
						Graph.Inputs = &profits[0]
					}),
					pieLink,
				)),
			),
		),
			container.NewTabItem("Numbers",
				container.NewVBox(
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
				)),
		))

	return scroll
}
