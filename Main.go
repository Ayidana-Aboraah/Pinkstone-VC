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

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject {
	titleText := widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewVBox(
		titleText,
		widget.NewButton("Sign-In", func() {
			var SignInStartUp dialog.Dialog
			var CreateUser dialog.Dialog

			nameEntry := widget.NewEntry()
			usrList := widget.NewListWithData(binding.BindStringList(&Database.Users), func() fyne.CanvasObject {
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
				Database.SaveData()
			}, w)

			SignInStartUp = dialog.NewConfirm("Sign In", "Create New User", func(b bool) {
				if b {
					CreateUser.Show()
				} else {
					dialog.ShowCustomConfirm("User Select", "Done", "Back", usrList, func(b bool) {
						if !b {
							SignInStartUp.Show()
						} else {
							titleText.SetText("Welcome " + Database.Users[Database.Current_User]) // change the title Text
						}
					}, w)
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
				text.SetText(Database.ItemKeys[val.ID].Name + " " + strconv.Itoa(int(val.Quantity)) + " by " + Database.Users[val.Usr])
			}

			reportList.OnSelected = func(id widget.ListItemID) {
				refunded_item = 0
			}

			dialog.ShowCustomConfirm("Refund", "Refund", "Cancel", reportList, func(b bool) {
				if !b {
					return
				}
				fmt.Println(refunded_item)
				if refunded_item < 0 {
					dialog.ShowInformation("Hmm", "No Refund selected, try selecting first", w)
					return
				}

				Database.RemoveReportEntry(0, refunded_item)
				UI.HandleError(Database.SaveData())

			}, w)

			//TODO: make a list indicating each sale in a report
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

						Database.Items = []Database.Item{}
						Database.ItemKeys = map[uint64]*Database.ItemEV{}
						Database.Reports = [2][]Database.Sale{}
						Database.Expenses = []Database.Expense{}
						Database.Free_Spaces = []int{}
						Database.Current_User = 0
						Database.Users = []string{""}

						label := w.Canvas().Content().(*fyne.Container).Objects[0].         //open vbox
						(*container.AppTabs).Items[0].Content.(*fyne.Container).Objects[0]. // Open Vbox in Main menu
						(*widget.Label)

						label.SetText("Welcome " + Database.Users[Database.Current_User]) // change the title Text
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
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil), widget.NewLabel(""))
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
		text.SetText(Database.ItemKeys[val.ID].Name + " x" + strconv.Itoa(int(val.Quantity)))

		btn.OnTapped = func() {
			shoppingCart = Database.DecreaseFromCart(val, shoppingCart)
			cartData.Set(Database.ConvertCart(shoppingCart))
			title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.ItemKeys[val.ID].Name + " x" + strconv.Itoa(int(val.Quantity)))
			shoppingList.Refresh()
		}
	}

	barcodeEntry := UI.NewNumEntry("Click and Scan")

	return container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(shoppingList),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b || len(shoppingCart) == 0 {
						return
					}
					receipt, total := Database.MakeReceipt(shoppingCart)
					shoppingCart = Database.BuyCart(shoppingCart)
					cartData.Set(Database.ConvertCart(shoppingCart))
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", total))
					// TODO: Get the data that would be added to the report
					// TODO: Send the data to the printer
					dialog.ShowInformation("Complete", receipt, w)
				}, w)
			}),

			widget.NewButton("Clear Cart", func() {
				cartData.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			}),

			widget.NewButton("New Item", func() {
				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", container.NewVBox(barcodeEntry), func(confirmed bool) {
					if !confirmed {
						return
					}

					id, err := strconv.ParseUint(barcodeEntry.Text, 10, 64)
					if err != nil {
						dialog.ShowInformation("Nope...", "Invalid Barcode", w)
						return
					}

					val, found := Database.ItemKeys[uint64(id)]
					if !found {
						dialog.ShowInformation("Oops", "Item not in database", w)
						return
					}

					dialog.ShowCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel("Is this the right item: "+val.Name)),
						func(b bool) {
							if !b {
								return
							}
							shoppingCart = Database.AddToCart(Database.ConvertItem(uint64(id)), shoppingCart)
							cartData.Set(Database.ConvertCart(shoppingCart))
							title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
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

	selectionEntry := UI.NewNumEntry("Year/Month")

	var profitDataSelect int
	var buttonType int

	dataSelectOptions := widget.NewSelect([]string{"Revenue", "Cost", "Profit"}, func(dataType string) {
		switch dataType {
		case "Revenue":
			profitDataSelect = 0
		case "Cost":
			profitDataSelect = 1
		case "Profit":
			profitDataSelect = 2
		}
	})

	reportDisplay := widget.NewLabel("")
	financeEntry := UI.NewNumEntry("YYYY/MM/DD [Select Custom again to select the date]")

	return container.NewVScroll(container.NewMax(container.NewVBox(
		widget.NewCard("Financial Reports", "", container.NewVBox(
			financeEntry,
			widget.NewSelect([]string{"Day", "Month", "Year", "Date"}, func(time string) {
				var variant uint8
				date := []uint8{}
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

					//String to Date conversion
					raw := strings.SplitN(financeEntry.Text, "/", 3)

					year, err := strconv.Atoi(raw[0][1:])
					if err != nil {
						fmt.Println("Something Seems up!") //DEBUG: REMOVE AFTER TESTING
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

				reportDisplay.SetText(Database.Report(variant, date))

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
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, 0)
				case 1:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, 1)
				case 2:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, profitDataSelect)
				case 3:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, 3)
				case 4:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 3, 0)
				}
			}),
			link,
		)),
	)))
}
