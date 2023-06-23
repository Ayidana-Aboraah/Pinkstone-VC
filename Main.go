package main

import (
	"BronzeHermes/Database"
	"BronzeHermes/Graph"
	"BronzeHermes/UI"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	go Graph.StartServer()

	Database.DataInit(false)

	// TODO: Access USB
	// TODO: Connect to Printer

	CreateWindow(a)
}

var w fyne.Window

func CreateWindow(a fyne.App) {
	w = a.NewWindow("Bronze Hermes")
	w.SetOnClosed(
		func() {
			Graph.StopSever()
			Database.SaveBackUp()
		},
	)

	if UI.HandleErrorWindow(Database.LoadData(), w) {
		dialog.ShowInformation("Back Up", "Loading BackUp", w)
		UI.HandleErrorWindow(Database.LoadBackUp(), w)
	}

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", makeMainMenu(a)),
		container.NewTabItem("Shop", makeShoppingMenu()),
		container.NewTabItem("Inventory", Database.MakeInfoMenu(w)),
		container.NewTabItem("Statistics", makeStatsMenu()),
	)))

	// Start Sign In Menu
	w.Content().(*fyne.Container).Objects[0].(*container.AppTabs).Items[0].Content.(*fyne.Container).Objects[1].(*widget.Button).OnTapped()

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject {
	var SignInStartUp dialog.Dialog
	var CreateUser dialog.Dialog

	usrData := binding.BindStringList(&Database.Users)
	titleText := widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewVBox(
		titleText,
		widget.NewButton("Sign-In", func() {

			nameEntry := widget.NewEntry()
			usrList := widget.NewListWithData(usrData, func() fyne.CanvasObject {
				return container.NewBorder(nil, nil, nil, nil, widget.NewLabel(""))
			}, func(di binding.DataItem, co fyne.CanvasObject) {
				text := co.(*fyne.Container).Objects[0].(*widget.Label)
				s, _ := di.(binding.String).Get()
				text.SetText(s)
			})

			usrList.OnSelected = func(id widget.ListItemID) {
				Database.Current_User = uint8(id)
			}

			CreateUser = dialog.NewForm("New", "Create", "Back", []*widget.FormItem{
				widget.NewFormItem("Username", nameEntry),
			}, func(b bool) {
				if !b {
					SignInStartUp.Show()
					return
				}

				Database.Users = append(Database.Users, nameEntry.Text)
				Database.Current_User = uint8(len(Database.Users) - 1)
				titleText.SetText("Welcome " + nameEntry.Text) // change the title Text
				usrData.Set(Database.Users)
				Database.SaveData()
			}, w)

			SignInStartUp = dialog.NewCustomConfirm("Sign In", "Login", "Create New", container.NewMax(usrList), func(b bool) {
				if b && len(Database.Users) > 0 {
					titleText.SetText("Welcome " + Database.Users[Database.Current_User])
				} else {
					CreateUser.Show()
				}
			}, w)

			SignInStartUp.Show()
		}),
		widget.NewButton("Refund", func() {
			refunded_item := -1

			reportData := binding.NewUntypedList()
			reportData.Set(Database.ConvertCart(Database.Reports[0]))
			reportList := widget.NewListWithData(reportData, func() fyne.CanvasObject {
				return container.NewBorder(nil, nil, nil, nil, widget.NewLabel(""))
			}, func(di binding.DataItem, co fyne.CanvasObject) {})

			reportList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
				text := obj.(*fyne.Container).Objects[0].(*widget.Label)
				v, _ := reportData.GetValue(id)
				val := v.(Database.Sale)
				text.SetText(Database.Item[val.ID].Name + " x" + strconv.Itoa(int(val.Quantity)) + " by " + Database.Users[val.Usr])
			}

			reportList.OnSelected = func(id widget.ListItemID) {
				refunded_item = id
			}

			dialog.ShowCustomConfirm("Refund", "Refund", "Cancel", reportList, func(b bool) {
				if !b {
					return
				}
				if refunded_item < 0 {
					dialog.ShowInformation("Hmm", "No Refund selected, try selecting first", w)
					return
				}

				Database.RemoveReportEntry(0, refunded_item)
				UI.HandleError(Database.SaveData())

			}, w)
		}),
		widget.NewButton("Save Backup Data", func() {
			go UI.HandleErrorWindow(Database.SaveBackUp(), w)
		}),
		widget.NewButton("Load Backup Data", func() {
			dialog.ShowInformation("Loading Back up Data", "Wait until back up is done loading...", w)
			go UI.HandleErrorWindow(Database.LoadBackUp(), w) // TODO: Add progress bar on main thread while waiting for this to happen
			dialog.ShowInformation("Loaded", "Back Up Loaded", w)
		}),
		widget.NewButton("Delete Database", func() {
			dialog.ShowConfirm("Are you sure?", "DELETE EVERYTHING",
				func(confirmed bool) {
					if !confirmed {
						return
					}
					dialog.ShowConfirm("You sure you sure?", "You sure you sure?", func(b bool) {
						if !confirmed {
							return
						}

						// Database.Items = []Database.Item{}
						Database.Item = map[uint64]*Database.Entry{}
						Database.Reports = [2][]Database.Sale{}
						Database.Expenses = []Database.Expense{}
						Database.Free_Spaces = []int{}
						Database.Current_User = 0
						Database.Users = []string{}
						CreateUser.Show()

						usrData.Set(nil)
						Database.InventoryData.Set(nil)
						Database.ExpenseData.Set(nil)
						titleText.SetText("Welcome")
						Database.SaveData()
					}, w)
				}, w)
		}),

		//Add inventory features here
	)
}

var shoppingCart []Database.Sale

func makeShoppingMenu() fyne.CanvasObject {

	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	cartData := binding.NewUntypedList()

	shoppingList := widget.NewListWithData(cartData,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("-", nil), widget.NewLabel(""))
		}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	shoppingList.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		cartData.Set(Database.ConvertCart(shoppingCart))
		title.SetText(fmt.Sprintf("Cart Total: %.2f", Database.GetCartTotal(shoppingCart)))
		shoppingList.Unselect(id)
	}

	shoppingList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		v, _ := cartData.GetValue(id)
		val := v.(Database.Sale)
		text.SetText(Database.Item[val.ID].Name + " x" + strconv.Itoa(int(val.Quantity)))

		btn.OnTapped = func() {
			shoppingCart = Database.DecreaseFromCart(val, shoppingCart)
			cartData.Set(Database.ConvertCart(shoppingCart))
			title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.Item[val.ID].Name + " x" + strconv.Itoa(int(val.Quantity)))
			shoppingList.Refresh()
		}
	}

	barcodeEntry := UI.NewNumEntry("Click and Scan")

	customerEntry := widget.NewEntry()

	return container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(shoppingList),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowForm("Do you want to buy all items in the Cart?", "Yes", "No",
					[]*widget.FormItem{widget.NewFormItem("Customer", customerEntry)}, func(b bool) {
						if !b || len(shoppingCart) == 0 {
							return
						}
						receipt := Database.MakeReceipt(shoppingCart, customerEntry.Text)
						shoppingCart = Database.BuyCart(shoppingCart)
						cartData.Set(Database.ConvertCart(shoppingCart))
						title.SetText("Cart Total: 0.0")

						// TODO: Get the data that would be added to the report
						// TODO: Send the data to the printer

						dialog.ShowCustom("Complete", "Done", container.NewVBox(
							widget.NewLabelWithStyle("PINKSTONE TRADING", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
							widget.NewLabelWithStyle(receipt, fyne.TextAlignCenter, fyne.TextStyle{}),
						), w)
					}, w)
			}),

			widget.NewButton("Clear Cart", func() {
				cartData.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			}),

			widget.NewButton("New Item", func() {
				barcodeEntry.SetText("")
				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", barcodeEntry, func(confirmed bool) {
					if !confirmed {
						return
					}

					id, err := strconv.ParseUint(barcodeEntry.Text, 10, 64)
					if err != nil {
						dialog.ShowInformation("Nope...", "Invalid Barcode", w)
						return
					}

					val, found := Database.Item[uint64(id)]
					if !found {
						dialog.ShowInformation("Oops", "Item not in database", w)
						return
					}

					barginEntry := UI.NewNumEntry("")

					dialog.ShowCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel(val.Name+"\nBargin Price?"), barginEntry),
						func(b bool) {
							if !b {
								return
							}
							s := Database.ConvertItem(uint64(id))
							if barginEntry.Text != "" {
								f, err := strconv.ParseFloat(barginEntry.Text, 32)
								if err != nil {
									dialog.ShowInformation("Invalid", "Invalid values passed into bargin entry", w)
								}
								s.Price = float32(f)
							}
							shoppingCart = Database.AddToCart(s, shoppingCart)
							cartData.Set(Database.ConvertCart(shoppingCart))
							title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
							shoppingList.Refresh()
						}, w)
				}, w)
			}),
		),
	)
}

func makeStatsMenu() fyne.CanvasObject {
	u, _ := url.Parse("http://localhost:8081/line")
	r, _ := url.Parse("http://localhost:8081/pie")

	link := widget.NewHyperlink("Go To Graph", u)

	selectionEntry := UI.NewNumEntry("YYYY-MM")

	var buttonType int

	reportDisplay := widget.NewLabel("")
	financeEntry := UI.NewNumEntry("YYYY-MM-DD")
	financeEntry.Hidden = true

	var variant uint8
	date := []uint8{}

	updateReport := func() {
		reportDisplay.SetText(Database.Report(variant, date))
	}

	updateGraph := func() {
		switch buttonType {
		case 0:
			Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 0, 0)
		case 1:
			Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, 1)
		case 2:
			Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 1, 0)
		}
	}

	content := container.NewVScroll(container.NewMax(container.NewVBox(
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
						reportDisplay.SetText("Type a date and select the date option again to get a report")
						return
					}

					raw := strings.SplitN(financeEntry.Text, "-", 3)

					year, err := strconv.Atoi(raw[0][1:])
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

					date = []uint8{uint8(day), uint8(month), uint8(year)}
				}

				updateReport()
			}),
			reportDisplay,
		)),
		widget.NewCard("Data Graphs", "", container.NewVBox(
			selectionEntry,
			widget.NewSelect([]string{"Items Graph", "Price Changes", "Item Popularity", "Item Sales"}, func(graph string) {
				switch graph {
				case "Items Graph":
					buttonType = 0
					link.URL = u
				case "Price Changes":
					buttonType = 1
					link.URL = u
				case "Item Popularity":
					buttonType = 2
					link.URL = r
				case "Item Sales":
					buttonType = 3
					link.URL = r
				case "Sales Over Time":
					buttonType = 4
					link.URL = u
				}
			}),
			link,
		)),
	)))

	content.OnScrolled = func(p fyne.Position) {
		updateReport()
		updateGraph()
	}

	return content
}
