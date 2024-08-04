package main

import (
	"BronzeHermes/Database"
	"BronzeHermes/Debug"
	"BronzeHermes/Graph"
	"BronzeHermes/Printer"
	"BronzeHermes/UI"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var printer Printer.Printer

func main() {
	a := app.NewWithID("PINKSTONE")
	// go Graph.StartServer()

	Database.DataInit()

	// TODO: Access USB
	// TODO: Connect to Printer

	CreateWindow(a)
}

var w fyne.Window

func CreateWindow(a fyne.App) {
	w = a.NewWindow("Pinkstone")
	w.SetOnClosed(
		func() {
			Graph.StopSever()
			Database.CleanUpDeadItems()
			Database.SaveData()
			Database.SaveBackUp()
			if printer.Active {
				printer.Close()
			}
		},
	)

	Debug.ShowError("Loading Data", Database.LoadData(), w)

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", mainMenu()),
		container.NewTabItem("Shop", shoppingMenu()),
		container.NewTabItem("Inventory", Database.ItemsMenu(w)),
		container.NewTabItem("Report", reportMenu()),
		// container.NewTabItem("Stats", makeStatsMenu()),
	)))

	w.Content().(*fyne.Container).Objects[0].(*container.AppTabs).OnSelected = func(ti *container.TabItem) {
		Database.RefreshInventory()
		RefreshReport()
		updateReport()
		// updateStatsGraphs()
	}

	w.ShowAndRun()
}

func mainMenu() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("System", func() {
			dialog.ShowCustom("System", "Done", container.NewVBox(
				widget.NewButton("Reset Sales", func() {
					dialog.ShowConfirm("Are You Sure", "Are You Sure you want to DELETE ALL THE SALES?", func(b bool) {
						if !b {
							return
						}
						Database.Sales = []Database.Sale{}
					}, w)
				}),
				widget.NewButton("Reset Items", func() {
					dialog.ShowConfirm("Are You Sure", "Are You Sure you want to DELETE ALL THE Items?", func(b bool) {
						if !b {
							return
						}
						Database.Items = map[uint16]*Database.Entry{}
					}, w)
				}),
			), w)
		}),
		widget.NewButton("Save Backup Data", func() {
			go Debug.ShowError("Backing up Data", Database.SaveBackUp(), w)
		}),

		widget.NewButton("Load Backup Data", func() {
			dialog.ShowConfirm("Are you Sure?", "Are you sure you want to load the backup data?", func(b bool) {
				if !b {
					return
				}
				Debug.ShowError("Loading Backup Data", Database.LoadBackUp(), w)
				dialog.ShowInformation("Loaded", "Back Up Loaded", w)
			}, w)
		}),

		widget.NewButton("Connect Printer", func() {
			devices, idxes, err := Printer.FindUSBDevices()
			Debug.ShowError("Looking for USB Devices", err, w)

			devListData := binding.NewIntList()
			devListData.Set(idxes)
			devList := widget.NewListWithData(devListData,
				func() fyne.CanvasObject {
					return container.NewBorder(nil, nil, nil, nil, widget.NewLabel(""))
				},
				func(di binding.DataItem, co fyne.CanvasObject) {
					idx, _ := di.(binding.Int).Get()
					co.(*fyne.Container).Objects[0].(*widget.Label).SetText(devices[idx].Product)
				})

			devList.OnSelected = func(id widget.ListItemID) {
				dialog.ShowCustomConfirm(devices[id].Product, "Select", "Close",
					widget.NewLabel(fmt.Sprintf("ID: %d\nProduct: %s\nVender ID: %d\nManufacturer: %s\nSerial: %s\nPath: %s",
						devices[id].ProductID, devices[id].Product, devices[id].VendorID, devices[id].Manufacturer, devices[id].Serial, devices[id].Path)), func(b bool) {

						if !b {
							return
						}

						p, e := Printer.Init(devices[id])

						if Debug.ShowError("Selecting "+devices[id].Product, e, w) {
							return
						}

						if printer.Active {
							printer.Close()
						}

						printer = p
					}, w)
				devList.Unselect(id)
			}

			dialog.ShowCustom("Devices", "Close", container.NewVBox(devList), w)
		}),

		widget.NewButton("Test Printing", func() {
			if printer.Active {
				entry := widget.NewEntry()
				dialog.ShowForm("", "Send", "Close", []*widget.FormItem{
					widget.NewFormItem("Message", entry),
				}, func(b bool) {
					if !b {
						return
					}
					printer.Print(Printer.ConvertForPrinter(entry.Text))
				}, w)
			} else {
				dialog.ShowInformation("", "No Printer Connected", w)
			}
		}),
	)
}

var shoppingCart []Database.Sale

func shoppingMenu() fyne.CanvasObject {

	title := widget.NewLabelWithStyle("Cart Total: 0.00", fyne.TextAlignCenter, fyne.TextStyle{})

	cartData := binding.NewUntypedList()

	shoppingList := widget.NewListWithData(cartData,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("-", nil), widget.NewLabel(""))
		}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	shoppingList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		val := shoppingCart[id]
		text.SetText(fmt.Sprintf("%s ₵%1.2f x%1.2f -> ₵%1.2f", Database.Items[val.ID].Name, val.Price, val.Quantity, val.Price*val.Quantity))
		btn.OnTapped = func() {
			shoppingCart = Database.DecreaseFromCart(id, shoppingCart)
			cartData.Set(Database.ConvertCart(shoppingCart))
			title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.Items[val.ID].Name + " x" + fmt.Sprint(val.Quantity))
			shoppingList.Refresh()
		}
	}

	shoppingList.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		cartData.Set(Database.ConvertCart(shoppingCart))
		title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
		shoppingList.Unselect(id)
	}

	customerEntry := UI.NewSearchBar("Customer Name Here...", Database.SearchCustomers)

	return container.NewVSplit(
		container.NewVSplit(
			title,
			container.NewStack(shoppingList),
		),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				customerEntry.SetText("")
				dialog.ShowForm("Do you want to buy all items in the Cart?", "Yes", "No",
					[]*widget.FormItem{widget.NewFormItem("Customer", customerEntry)}, func(b bool) {
						if !b || len(shoppingCart) == 0 {
							return
						}

						i := customerEntry.ResultOrCreate(func() {
							Database.Customers = append(Database.Customers, customerEntry.Text)
						})

						receipt := Database.MakeReceipt(shoppingCart, customerEntry.Text)
						shoppingCart = Database.BuyCart(shoppingCart, uint16(i))
						cartData.Set(Database.ConvertCart(shoppingCart))

						title.SetText("Cart Total: 0.00")

						Debug.ShowError("Saving Data", Database.SaveData(), w)

						dialog.ShowCustomConfirm("Complete", "Print", "Done", container.NewVBox(
							widget.NewLabelWithStyle("PINKSTONE TRADING", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
							widget.NewLabelWithStyle(receipt, fyne.TextAlignCenter, fyne.TextStyle{}),
						), func(printing bool) {
							if printing && printer.Active {
								printer.Print(Printer.ConvertForPrinter(receipt))
							}
						}, w)
					}, w)
			}),

			widget.NewButton("Clear Cart", func() {
				cartData.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			}),

			widget.NewButton("New Item", func() {
				searchBar := UI.NewSearchBar("Item Name Here...", Database.SearchInventory)

				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", searchBar, func(confirmed bool) {
					if !confirmed {
						return
					}

					id := searchBar.Result()

					if Debug.HandleKnownError(0, id < 0, w) {
						return
					}

					val := Database.Items[uint16(id)]

					barginEntry := UI.NewNumEntry("A bargined price for the customer")
					quantityEntry := UI.NewNumEntry("The amount being bought")

					options := widget.NewAccordion(
						widget.NewAccordionItem("Bargin Price", barginEntry),
						widget.NewAccordionItem("Pieces", quantityEntry))

					var menu dialog.Dialog
					menu = dialog.NewCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel(val.Name), options),
						func(b bool) {
							if !b {
								return
							}
							s := Database.NewItem(uint16(id))

							id := Database.ProcessNewItemData(barginEntry.Text, quantityEntry.Text, &s)
							if Debug.HandleKnownError(id, id != -1, w) {
								menu.Show()
							}

							ShoppingCart, errID := Database.AddToCart(s, shoppingCart)
							Debug.HandleKnownError(errID, errID != Debug.Success, w)

							shoppingCart = ShoppingCart
							cartData.Set(Database.ConvertCart(shoppingCart))
							title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))

							shoppingList.Refresh()
							shoppingList.ScrollToBottom()
						}, w)
					menu.Show()
				}, w)
			}),
		),
	)
}

var updateReport func()

var RefreshReport func()

func reportMenu() fyne.CanvasObject {

	reportDisplay := widget.NewLabel("")
	financeEntry := UI.NewNumEntry("YYYY-MM-DD")
	financeEntry.Hidden = true

	var variant uint8
	date := []int{}

	updateReport = func() {
		str, _ := Database.CompileReport(variant, date)
		reportDisplay.SetText(str)
	}

	customerSearch := UI.NewSearchBar("Customer Name Here...", Database.SearchCustomers)

	reportData := binding.NewUntypedList()
	reportData.Set(Database.ConvertCart(Database.Sales))

	RefreshReport = func() {
		reportData.Set(Database.ConvertCart(Database.Sales))
	}

	reportList := widget.NewListWithData(reportData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel(""))
	}, func(di binding.DataItem, co fyne.CanvasObject) {
		v, _ := di.(binding.Untyped).Get()
		val := v.(Database.Sale)
		display := co.(*fyne.Container).Objects[0].(*widget.Label)
		out := fmt.Sprintf("%s x%1.2f for ₵%1.2f Customer: %s, [%s]",
			Database.Items[val.ID].Name, val.Quantity, val.Price*val.Quantity, Database.Customers[val.Customer], time.Unix(val.Timestamp, 0).Local().String())

		if val.Customer == 0 {
			out = "[DAMAGED] " + out
		}

		display.SetText(out)

	})

	reportList.OnSelected = func(id widget.ListItemID) {
		v, _ := reportData.GetValue(id) // Maybe Log
		val := v.(Database.Sale)
		var infoText string

		infoText = fmt.Sprintf("Name: %s\nPrice: %1.2f\nCost: %1.2f\nQuantity: %1.2f\nTotal Revenue: %1.2f\nTotal Profit: %1.2f\nCustomer: %s",
			Database.Items[val.ID].Name, val.Price, val.Cost, val.Quantity, val.Price*val.Quantity, (val.Price-val.Cost)*val.Quantity, Database.Customers[val.Customer])

		if val.Customer == 0 {
			infoText = "[DAMAGED] " + infoText
		}

		dialog.ShowCustomConfirm("Info", "Refund", "Close", widget.NewLabel(infoText), func(b bool) {
			if !b {
				return
			}

			Database.RemoveFromSales(id)
			Debug.ShowError("Saving Data", Database.SaveData(), w)
			reportData.Set(Database.ConvertCart(Database.Sales))
			updateReport()
			reportList.Refresh()
			fmt.Printf("Refunded item")
		}, w)
		reportList.Unselect(id)
	}

	content := container.New(layout.NewGridLayoutWithRows(3),
		widget.NewCard("Financial Reports", "", container.NewVBox(
			financeEntry,
			widget.NewSelect([]string{"Day", "Month", "Year", "Date"}, func(time string) {
				financeEntry.Hidden = true

				switch time {
				case "Day":
					variant = Database.ONCE
				case "Month":
					variant = Database.MONTHLY
				case "Year":
					variant = Database.YEARLY
				case "Date": //The user will have to double tap when using Dates
					financeEntry.Hidden = false
					if financeEntry.Text == "" {
						reportDisplay.SetText("Type a date and selefct the date option again to get a report")
						return
					}

					raw := strings.SplitN(financeEntry.Text, "-", 3)

					year, err := strconv.Atoi(raw[0])
					if err != nil {
						return
					}

					variant = Database.YEARLY

					var month int
					var day int

					if len(raw) > 1 {
						month, _ = strconv.Atoi(raw[1]) // NOTE: Error handling may be needed here, unknown for now
						variant = Database.MONTHLY
					}

					if len(raw) > 2 {
						day, _ = strconv.Atoi(raw[2])
						variant = Database.ONCE
					}

					date = []int{day, month, year}
				}

				updateReport()
			}),
			reportDisplay,
		)),

		widget.NewCard("Sales", "", container.NewVBox(
			customerSearch,
			widget.NewButton("Search", func() {
				customerIdx := customerSearch.Result()
				found := []Database.Sale{}

				if customerIdx == -1 {
					found = Database.Sales
				} else {
					for _, v := range Database.Sales {
						if v.Customer == uint16(customerIdx) {
							found = append(found, v)
						}
					}
				}

				reportData.Set(Database.ConvertCart(found))
				reportList.Refresh()
			}),
		)),
		container.NewStack(reportList),
	)
	return content
}

var updateStatsGraphs func()

// func makeStatsMenu() fyne.CanvasObject {
// 	dateEntry := UI.NewNumEntry("YYYY-MM-DD")
// 	itemSearch := UI.NewSearchBar("Item Name Here...", Database.SearchInventory)
// 	customerSearch := UI.NewSearchBar("Customer Name Here...", Database.SearchCustomers)

// 	currentGraphType := 0
// 	graphSelect := widget.NewSelect([]string{"Over Time", " Current Total"}, func(s string) {
// 		switch s {
// 		case "Over Time":
// 			currentGraphType = 0
// 		case "Current Total":
// 			currentGraphType = 1
// 		}
// 		updateStatsGraphs()
// 	})

// 	dataType := 0
// 	dataTypeSelect := widget.NewSelect([]string{"Profit", "Cost", "Revenue"}, func(s string) {
// 		switch s {
// 		case "Profit":
// 			dataType = 0
// 		case "Cost":
// 			dataType = 1
// 		case "Revenue":
// 			dataType = 2
// 		}
// 		updateStatsGraphs()
// 	})
// 	dateEntry.OnChanged = func(s string) {
// 		updateStatsGraphs()
// 	}
// 	customerSearch.OnChanged = func(s string) {
// 		updateStatsGraphs()
// 	}
// 	itemSearch.OnChanged = func(s string) {
// 		updateStatsGraphs()
// 	}

// 	updateStatsGraphs = func() {
// 		// errID := -1

// 		switch currentGraphType {
// 		case 0:
// 			Graph.Labels, Graph.LineInputs, _ = Database.GetLine(dateEntry.Text, dataType, itemSearch.Result(), customerSearch.Result())
// 			// if Debug.HandleKnownError(errID, errID != Debug.Success, w) {
// 			// 	return
// 			// }
// 		case 1:
// 			Graph.Labels, Graph.Inputs, _ = Database.GetPie(dateEntry.Text, dataType, itemSearch.Result(), customerSearch.Result())
// 			// if Debug.HandleKnownError(errID, errID != Debug.Success, w) {
// 			// 	return
// 			// }
// 		}
// 	}

// 	lineURL, _ := url.ParseRequestURI("localhost:8081/line")
// 	pieURL, _ := url.ParseRequestURI("localhost:8081/pie")

// 	lineLink := widget.NewHyperlinkWithStyle("Line Graph", lineURL, fyne.TextAlignCenter, fyne.TextStyle{})
// 	pieLink := widget.NewHyperlinkWithStyle("Pie Graph", pieURL, fyne.TextAlignCenter, fyne.TextStyle{})

// 	return container.NewVBox(
// 		widget.NewCard("Item Sales", "", container.NewVBox(
// 			graphSelect,
// 			dataTypeSelect,
// 			dateEntry,
// 			itemSearch,
// 			customerSearch,
// 			lineLink,
// 			pieLink,
// 		)),
// 	)
// }
